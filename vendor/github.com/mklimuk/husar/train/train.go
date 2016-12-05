package train

import (
	"fmt"
	"time"

	"github.com/mklimuk/husar/db"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
)

// ServiceType represents train services
type ServiceType string

// Services catalog
const (
	REZO ServiceType = "REZO"
	KOND ServiceType = "KOND"
)

// AnnonOption represents announcements-specific train options
type AnnonOption string

// Predefined options
const (
	ArrivalText    AnnonOption = "arrivalText"
	DepartureText  AnnonOption = "departureText"
	DelayText      AnnonOption = "delayText"
	SetupText      AnnonOption = "setupText"
	BusSetupText   AnnonOption = "busSetupText"
	BusArrivalText AnnonOption = "busArrivalText"
	Pause          AnnonOption = "pause"
	I18n           AnnonOption = "i18n"
)

// Train represents rethink db Trains table schema
type Train struct {
	Arrival        *TimetableEvent `gorethink:"arrival" json:"arrival"`
	Departure      *TimetableEvent `gorethink:"departure" json:"departure"`
	FirstEvent     *time.Time      `gorethink:"firstEvent" json:"firstEvent,omitempty"`
	FirstLiveEvent *time.Time      `gorethink:"firstLiveEvent,omitempty" json:"firstLiveEvent,omitempty"`
	Delay          int             `gorethink:"delay" json:"delay"`
	Route          *Route          `gorethink:"route" json:"route"`
	Carrier        string          `gorethink:"carrier" json:"carrier"`
	Category       string          `gorethink:"category" json:"category"`
	Day            time.Time       `gorethink:"day" json:"day"`
	ID             string          `gorethink:"id,omitempty" json:"id"`
	OrderID        int             `gorethink:"orderId" json:"orderId"`
	StationID      int             `gorethink:"stationId" json:"stationId"`
	StationName    string          `gorethink:"stationName" json:"stationName"`
	TrainID        string          `gorethink:"trainId" json:"trainId"`
	Name           string          `gorethink:"name" json:"name"`
	Services       *[]Service      `gorethink:"services,omitempty" json:"services,omitempty"`
	Settings       *Settings       `gorethink:"options,omitempty" json:"options,omitempty"`
}

func (t *Train) String() string {
	return t.ID
}

// TimetableEvent holds all temporal data about given train (it is used for arrival and departure events).
// ExpectedTime is added for client applications convienience so that it doesn't have to calculate anything.
type TimetableEvent struct {
	Time         time.Time `gorethink:"time" json:"time"`
	ExpectedTime time.Time `json:"expectedTime"`
	Platform     string    `gorethink:"platform" json:"platform"`
	Track        int       `gorethink:"track" json:"track"`
}

// Station represents one station with optional arrival and departure
type Station struct {
	Arrival       *time.Time `gorethink:"arrival,omitempty" json:"arrival,omitempty"`
	Departure     *time.Time `gorethink:"departure,omitempty" json:"departure,omitempty"`
	ID            int        `gorethink:"id" json:"id"`
	Name          string     `gorethink:"name" json:"name"`
	Important     bool       `gorethink:"important" json:"important"`
	NumberOnRoute int        `gorethink:"numberOnRoute" json:"numberOnRoute"`
}

// Service represents services available onboard
type Service struct {
	ID       string   `gorethink:"id" json:"id"`
	Name     string   `gorethink:"name" json:"name"`
	Carriage []string `gorethink:"carriage,omitempty" json:"carriage,omitempty"`
}

// Route represents all information about route of train
type Route struct {
	CurrentStationPositionOnRoute int       `gorethink:"currentStationPositionOnRoute" json:"currentStationPositionOnRoute"`
	CurrentStationOnSubroute      bool      `gorethink:"currentStationOnSubroute" json:"currentStationOnSubroute"`
	EndStation                    Station   `gorethink:"endStation" json:"endStation"`
	StartStation                  Station   `gorethink:"startStation" json:"startStation"`
	SubrouteStart                 *Station  `gorethink:"subrouteStart,omitempty" json:"subrouteStart,omitempty"`
	SubrouteEnd                   *Station  `gorethink:"subrouteEnd,omitempty" json:"subrouteEnd,omitempty"`
	Stations                      []Station `gorethink:"stations" json:"stations"`
}

/*
Store methods related to Trains
*/

func (store *rethinkStore) Get(id string) (*Train, error) {
	resp, err := r.DB(store.db).Table(store.trainTable).Get(id).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "Get"}).
			WithError(err).Error("Error while reading train.")
		return nil, err
	}
	if resp.IsNil() {
		return nil, err
	}
	a := new(Train)
	err = resp.One(a)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "Get"}).
			WithError(err).Error("Error while parsing train.")
		return nil, err
	}
	return a, err
}

func (store *rethinkStore) Save(t *Train) (*Train, error) {
	resp, err := r.DB(store.db).Table(store.trainTable).Insert(t, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "Save", "train": t.ID}).
			WithError(err).Error("Error occured while saving train.")
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "Save", "train": t.ID, "result": resp}).
			Debug("Database save result.")
	}
	return t, err
}

func (store *rethinkStore) SaveAll(tr *[]Train) error {
	res, err := r.DB(store.db).Table(store.trainTable).Insert(*tr, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "SaveAll", "trains": *tr}).
			WithError(err).Error("Error occured while saving train objects.")
		return err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "SaveAll", "trains": *tr, "result": res}).
			Debugln("Database save result.")
	}
	return err
}

func (store *rethinkStore) GetBetween(start, end *time.Time) (trains []Train, err error) {
	var res *r.Cursor
	/*r.DB(store.db).Table(store.trainTable).
		Between(start, end, r.BetweenOpts{Index: db.TrainTimeIndex}).OuterJoin(r.DB(store.db).Table(store.realtimeTable), func(t r.Term, re r.Term) r.Term {
		return t.Field("id").Eq(re.Field("id"))
	}).Zip().Map(r.Branch(r.Row.HasFields("firstLiveEvent"), r.Row, r.Row.Merge(func(t r.Term) map[string]interface{} {
		return map[string]interface{}{
			"firstLiveEvent": t.Field("firstEvent"),
		}
	}))).OrderBy(db.RealtimeTimeIndex).Run(store.session)*/
	if res, err = r.DB(store.db).Table(store.trainTable).
		Between(start, end, r.BetweenOpts{Index: db.TrainTimeIndex}).
		Map(func(t r.Term) r.Term {
			return t.Merge(r.DB(store.db).Table(store.realtimeTable).Get(t.Field("id")).Default(map[string]interface{}{
				"firstLiveEvent": t.Field("firstEvent"),
			}))
		}).OrderBy(db.RealtimeTimeIndex).Run(store.session); err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "GetBetween", "start": start, "end": end}).
			WithError(err).Error("Error executing query on the timetable table.")
		return
	}
	if err = res.All(&trains); err != nil {
		log.WithFields(log.Fields{"logger": "train.store.train", "method": "GetBetween", "start": start, "end": end}).
			WithError(err).Error("Error parsing results from the realtime table.")
		return
	}
	return
}

func (store *rethinkStore) GetBetweenWithRealtime(start, end *time.Time) (res []Train, IDs []string, err error) {

	clog := log.WithFields(log.Fields{"logger": "shared.model.train", "method": "GetBetweenWithRealtime"})
	// get regular trains for period
	var regular []Train
	if regular, err = (*store).GetBetween(start, end); err != nil {
		clog.WithError(err).Error("Could not load trains.")
		return res, IDs, err
	}
	// get live trains for period
	var live []Train
	if live, err = (*store).GetRealtimeBetween(start, end); err != nil {
		clog.WithError(err).Error("Could not load live trains.")
		return res, IDs, err
	}
	// merge the two lists
	res = append(regular, live...)
	res, IDs = removeDuplicates(&res, clog)
	return res, IDs, err
}

func (store *rethinkStore) Search(query string) (t []*Train, err error) {
	clog := log.WithFields(log.Fields{"logger": "shared.model.train", "method": "Search"})
	match := fmt.Sprintf("(?i).*%s.*", query) // case insensitive
	var cur *r.Cursor
	if cur, err = r.DB(store.db).Table(store.trainTable).Filter(func(t r.Term) r.Term {
		return t.Field("name").Match(match).Or(t.Field("trainId").Match(match)).
			Or(t.Field("route").Field("startStation").Field("name").Match(match)).
			Or(t.Field("route").Field("endStation").Field("name").Match(match))
	}, r.FilterOpts{}).Run(store.session); err != nil {
		clog.WithError(err).Error("Error getting search results.")
		return
	}
	defer cur.Close()
	if cur.IsNil() {
		clog.Info("No search results returned.")
		return
	}
	cur.All(&t)
	return
}

func removeDuplicates(list *[]Train, clog *log.Entry) (result []Train, IDs []string) {

	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("size", len(*list)).Debug("Removing duplicates from trains list.")
	}

	// build a hashmap to eliminate duplicates
	var distinct = make(map[string]*Train, len(*list))
	for _, t := range *list {
		if log.GetLevel() >= log.DebugLevel {
			clog.WithField("train", t.ID).Debug("Checking train.")
		}

		if _, ok := distinct[t.ID]; !ok {

			if log.GetLevel() >= log.DebugLevel {
				clog.WithField("train", t.ID).Debug("Train not yet added. Adding.")
			}

			distinct[t.ID] = &t
			result = append(result, t)
			IDs = append(IDs, t.ID)
		}
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithField("size", len(result)).Debug("Duplicates removed.")
	}
	return
}
