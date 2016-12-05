package service

import (
	"encoding/json"
	"time"

	"github.com/mklimuk/husar/event"
	"github.com/mklimuk/husar/msg"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

//Timetable is a timetable service
type Timetable interface {
	SearchTrain(query string, removeDuplicates bool) (t []*train.Train, err error)
}

type timetable struct {
	b     *event.Bus
	store train.Store
}

//NewTimetable it the timetable service constructor
func NewTimetable(b *event.Bus, store train.Store) Timetable {
	s := timetable{b: b, store: store}
	if b != nil {
		b.Subscribe(event.GetTimetable, s.getTimetable)
		b.Subscribe(event.TrainUpdateEvent, s.updateTrain)
		b.Subscribe(event.TrainUpdateMode, s.updateTrain)
		b.Subscribe(event.AudioSettingsUpdate, s.updateSettings)
	}
	return Timetable(&s)
}

func (s *timetable) getTimetable(request string) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable", "request": request}).
			Debug("Received get timetable request.")
	}
	var err error
	var req = new(msg.TimetableRequest)
	if err = json.Unmarshal([]byte(request), req); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable"}).
			WithError(err).Error("Error unmarshalling request")
		return
	}

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable", "from": req.From, "to": req.To}).
			Debug("Parsed parameters.")
	}
	var timetable []train.Train
	var IDs []string
	if timetable, IDs, err = s.store.GetBetweenWithRealtime(&req.From, &req.To); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable"}).
			WithError(err).Error("Error retrieving timetable contents")
		return
	}
	var settings *[]train.Settings
	if len(IDs) > 0 {
		if settings, err = s.store.GetAllSettings(IDs...); err != nil {
			log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable"}).
				WithError(err).Error("Error retrieving train settings")
		}
	}

	result := msg.TimetableResponse{StationID: req.StationID, Timetable: []*train.Train{}}
	set := train.BuildSettingsMap(settings)
	for _, t := range timetable {
		if s := set[t.ID]; s.ID == t.ID {
			t.Settings = &s
		}
		result.Timetable = append(result.Timetable, processTrain(t))
	}

	var res []byte
	if res, err = json.Marshal(result); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "getTimetable"}).
			WithError(err).Error("Error marshalling response")
	}
	go s.b.Publish(event.TimetableContent, event.TimetableContent, string(res))
}

func (s *timetable) updateTrain(req string) {

	clog := log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "updateTrain"})
	var m *msg.TrainUpdate
	var err error
	if m, err = msg.ParseTrainUpdate(req); m == nil {
		clog.WithError(err).Error("Could not parse message")
		return
	}
	var t *train.Train
	if t, err = s.store.Get(m.TrainID); err != nil {
		clog.WithError(err).Error("Could not load train from the database")
		return
	}
	// we get settings for this train
	var set *[]train.Settings
	if set, err = s.store.GetAllSettings(m.TrainID); err != nil {
		clog.WithError(err).Error("Could not load train from the database")
		return
	}
	if len(*set) > 0 {
		t.Settings = &(*set)[0]
	}

	if t.Settings == nil {
		initSettings(t)
	}
	clog.WithField("settings", t.Settings.ID).Info("Updating settings")
	if m.Mode != nil {
		t.Settings.Mode = (*m.Mode)
	}
	if m.Event != nil {
		if m.Event.Delay != nil {
			updateDelay(t, m.Event.Delay)
		}
		switch m.Event.Type {
		case msg.Arrival:
			updateEvent(t.Settings.Arrival, m.Event)
		case msg.Departure:
			updateEvent(t.Settings.Departure, m.Event)
		case msg.All:
			updateEvent(t.Settings.Arrival, m.Event)
			updateEvent(t.Settings.Departure, m.Event)
		}
	}
	clog.WithField("settings", t.Settings.ID).Info("Saving settings")
	if _, err = s.store.SaveSettings(t.Settings); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "updateTrain"}).
			WithError(err).Error("Error saving settings to the database")
	}
	var res []byte
	if res, err = json.Marshal(processTrain(*t)); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "updateTrain"}).
			WithError(err).Error("Could not parse update message")
	}
	go s.b.Publish(event.TrainUpdate, event.TrainUpdate, string(res))
}

func initSettings(t *train.Train) {
	t.Settings = &train.Settings{Type: train.TrainInstance}
	if t.Arrival != nil {
		t.Settings.Arrival = new(train.TimetableEvent)
		*t.Settings.Arrival = *t.Arrival
		t.Settings.Arrival.Time = t.Arrival.Time
	}
	if t.Departure != nil {
		t.Settings.Departure = new(train.TimetableEvent)
		*t.Settings.Departure = *t.Departure
		t.Settings.Departure.Time = t.Departure.Time
	}
	t.Settings.Mode = train.Auto
	t.Settings.Delay = t.Delay
	t.Settings.ID = t.ID
}

func updateDelay(t *train.Train, delay *int) {
	t.Settings.Mode = train.Manual
	t.Settings.Delay = *delay
}

func updateEvent(ev *train.TimetableEvent, e *msg.Event) {
	if ev == nil {
		return
	}
	if e.Platform != nil {
		ev.Platform = *e.Platform
	}
	if e.Track != nil {
		ev.Track = *e.Track
	}
}

// Processing train is there to simplify front end logic. We apply all required
// modifications so that train structure sent to the client is as simple as possible
func processTrain(t train.Train) *train.Train {
	clog := log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "processTrain", "train": t.ID})
	t.Delay = train.GetDelay(&t)
	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("delay", t.Delay).Debug("Processing train for client use")
	}
	t.Category = train.GetCategory(&t)
	t.Arrival = train.GetArrival(&t)
	t.Departure = train.GetDeparture(&t)
	if t.Arrival != nil {
		t.Arrival.ExpectedTime = t.Arrival.Time.Add(time.Duration(t.Delay) * time.Minute)
		if log.GetLevel() >= log.DebugLevel {
			clog.WithField("arrival", t.Arrival.ExpectedTime).Debug("Expected arrival")
		}
	}
	if t.Departure != nil {
		t.Departure.ExpectedTime = t.Departure.Time.Add(time.Duration(t.Delay) * time.Minute)
		if log.GetLevel() >= log.DebugLevel {
			clog.WithField("arrival", t.Departure.ExpectedTime).Debug("Expected departure")
		}
	}
	t.FirstEvent = nil
	t.FirstLiveEvent = nil
	return &t
}
