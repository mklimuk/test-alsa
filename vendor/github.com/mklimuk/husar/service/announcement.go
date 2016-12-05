package service

import (
	e "errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/audio"
	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/errors"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

//Announcement watches for train table changes and synchronizes announcements accordingly
type Announcement interface {
	ProcessRealtimeChange(next, old *train.Realtime)
	ProcessSettingsChange(next, old *train.Settings)
	GeneratePreview(ID string, params annon.TemplateParams) (txt *annon.Text, length int, file string, err error)
	GenerateFromCatalog(ID string, zoneID string, startTime *time.Time, endTime *time.Time, dayInterval int, nightInterval int, params annon.TemplateParams) (err error)
	GenerateForPeriod(start, end *time.Time)
	GetCatalog() (catalog map[string][]*annon.CatalogTemplate, err error)
	GetTemplate(templateID string) (*annon.CatalogTemplate, error)
	SaveTemplate(templ *annon.CatalogTemplate) (*annon.CatalogTemplate, error)
	DeleteTemplate(templateID string) error
}

type announcement struct {
	gen    *annon.Generator
	cat    *annon.Catalog
	store  *annon.Store
	trains *train.Store
	audio  *audio.Catalog
}

// NewAnnouncement Service implementation constructor
func NewAnnouncement(gen *annon.Generator, store *annon.Store, trains *train.Store, cat *annon.Catalog, au *audio.Catalog) Announcement {
	s := &announcement{
		gen:    gen,
		store:  store,
		trains: trains,
		cat:    cat,
		audio:  au,
	}
	return Announcement(s)
}

func (s *announcement) ProcessRealtimeChange(next, old *train.Realtime) {

	// context logger
	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "ProcessRealtimeChange"})

	// handle delete
	if next == nil {
		if log.GetLevel() >= log.InfoLevel {
			clog.Info("Received a delete update. Ignoring.")
		}
		return
	}

	clog = clog.WithField("realtime", next.ID)
	if log.GetLevel() >= log.InfoLevel {
		clog.WithField("delay", next.Delay).Info("The delay has changed. Processing changes.")
	}

	var t *train.Train
	var err error
	if t, err = (*s.trains).Get(next.ID); err != nil {
		clog.Error("Error while retrieving train from the store. Cannot process realtime change.")
		return
	}
	if t == nil {
		// we didn't find the corresponding train which is weird
		clog.Warn("Train not found corresponding to the realtime info id.")
		return
	}
	prev := new(train.Train)
	// we clone the train to be able to provide
	// both previous and current train to providers
	train.Copy(t, prev)
	t.Delay = next.Delay
	t.FirstLiveEvent = &next.FirstEvent
	if old != nil {
		prev.Delay = old.Delay
		prev.FirstLiveEvent = &old.FirstEvent
	}
	settings := s.GetSettingsForTrain(t.ID, t.TrainID, t.OrderID)
	t.Settings = settings
	prev.Settings = settings
	s.ProcessTrain(t, prev)

}

func (s *announcement) GetSettingsForTrain(ID string, trainID string, orderID int) *train.Settings {
	// get all available Settings
	// concatenate intelligently
	var settings *[]train.Settings
	var err error
	if settings, err = (*s.trains).GetAllSettings(ID, trainID, strconv.Itoa(orderID)); err != nil {
		log.WithFields(log.Fields{"logger": "annon.service", "method": "GetSettingsForTrain", "ID": ID, "trainID": trainID, "orderID": orderID}).
			WithError(err).Error("Could not read settings")
	}
	if len(*settings) > 0 {
		return &(*settings)[0]
	}
	return nil
}

func (s *announcement) ProcessSettingsChange(next, old *train.Settings) {
	// TODO think of a mechanism that we should use to update multiple trains
	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "ProcessSettingsChange"})
	if next == nil {
		clog.Warn("Removing settings. This should never happen. Assuming it is a maintenance operation.")
		return
	}
	var t *train.Train
	var err error

	if t, err = (*s.trains).Get(next.ID); err != nil {
		clog.WithField("settings", next.ID).WithError(err).Error("Error getting the train associated with settings")
		return
	}
	if t == nil {
		clog.WithField("settings", next.ID).Warn("Train not found for settings. Ignoring.")
		return
	}

	var r *train.Realtime
	if r, err = (*s.trains).GetRealtime(next.ID); err != nil {
		clog.WithField("settings", next.ID).WithError(err).Error("Error getting realtime for train")
		return
	}

	prev := new(train.Train)
	// we clone the train to be able to provide
	// both previous and current train to providers
	train.Copy(t, prev)
	if r != nil {
		prev.Delay = r.Delay
		t.Delay = r.Delay
	}
	t.Settings = next
	prev.Settings = old
	s.ProcessTrain(t, prev)

}

func (s *announcement) GenerateForPeriod(start, end *time.Time) {

	// context logger
	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "GenerateForPeriod"})

	if start == nil || end == nil {
		clog.WithFields(log.Fields{"start": start, "end": end}).
			Error("Invalid parameters provided to method. Ignoring.")
	}

	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"start": *start, "end": *end}).
			Debug("Generating announcements for period.")
	}

	timetable, IDs, err := (*s.trains).GetBetweenWithRealtime(start, end)
	if err != nil {
		clog.WithError(err).Error("Error getting realtime trains.")
		return
	}

	if len(timetable) == 0 {
		if log.GetLevel() >= log.InfoLevel {
			clog.WithFields(log.Fields{"start": *start, "end": *end}).
				Info("No trains to process in the requested period. Exiting.")
		}
		return
	}

	// get settings for period
	var set *[]train.Settings
	if set, err = (*s.trains).GetAllSettings(IDs...); err != nil {
		clog.WithError(err).Error("Could not load settings for period.")
		return
	}
	settings := train.BuildSettingsMap(set)

	//processing trains
	var ts train.Settings
	for _, t := range timetable {
		ts = settings[t.ID]
		t.Settings = &ts
		s.ProcessTrain(&t, nil)
	}
}

func (s *announcement) ProcessTrain(train, old *train.Train) error {

	// context logger
	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "ProcessTrain"})
	if log.GetLevel() >= log.DebugLevel {
		clog.Debug("Processing train.")
	}

	if train == nil {
		log.Warn("Received a nil train for processing. Trying to delete old announcements if possible.")
		if old != nil {
			if err := (*s.store).DeleteForTrain(old.ID); err != nil {
				clog.WithField("train", old.ID).WithError(err).
					Error("Error trying to cleanup.")
			}
		}
		return nil
	}
	clog = clog.WithField("train", train.ID)

	currentAnnons, err := (*s.store).GetForTrain(train.ID)
	if err != nil {
		clog.WithError(err).Error("Error while getting current announcements.")
		return err
	}
	now := time.Now()
	annons, toDelete, err := (*s.gen).ForTrain(train, old, currentAnnons, &now, s.cat)
	if err != nil {
		clog.WithError(err).Error("Error while generating announcements. Aborting treatment.")
		return err
	}
	for _, a := range toDelete {
		if log.GetLevel() >= log.DebugLevel {
			clog.WithField("annon", a.ID).Debug("Deleting announcement.")
		}
		(*s.store).Delete(a.ID)
	}
	var id string
	var duration int
	var approxLen bool
	for _, a := range annons {
		if log.GetLevel() >= log.DebugLevel {
			clog.WithField("annon", a.ID).Debug("Generating announcement audio.")
		}
		if id, _, duration, approxLen, err = (*s.audio).Generate(a.Text.TtsText); err != nil {
			clog.WithField("annon", a.ID).WithError(err).Error("Could not generate audio file.")
			continue
		}
		aud := annon.Audio{Duration: duration, FileID: id, ApproxLen: approxLen}
		a.Audio = &aud
	}
	if err = (*s.store).SaveAll(&annons); err != nil {
		clog.WithError(err).Errorln("Error saving new announcements.")
		return err
	}
	return nil
}

func (s *announcement) GeneratePreview(catalogID string, params annon.TemplateParams) (txt *annon.Text, duration int, file string, err error) {

	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "GeneratePreview"})
	var templ *annon.CatalogTemplate
	if templ, err = (*s.cat).Get(catalogID); err != nil {
		clog.WithError(err).Error("Could not retrieve catalog template.")
		return
	}
	if templ == nil {
		clog.Debug("Template not found")
		return nil, -1, "", errors.NewError("Template not found.", errors.NotFound)
	}
	clog.Debug("Loaded requested template.")
	if txt, err = (*s.gen).Preview(templ, params); err != nil {
		clog.WithError(err).Error("Could not generate announcements for template.")
		return
	}

	if file, _, duration, _, err = (*s.audio).Generate(txt.TtsText); err != nil {
		clog.WithError(err).Error("Could not generate audio file.")
		return
	}
	return
}

func (s *announcement) GenerateFromCatalog(catalogID string, zoneID string, startTime *time.Time, endTime *time.Time, dayInterval int, nightInterval int, params annon.TemplateParams) (err error) {

	clog := log.WithFields(log.Fields{"logger": "annongen.service", "method": "GenerateFromCatalog"})

	dayInt := time.Duration(dayInterval) * time.Minute
	nightInt := time.Duration(nightInterval) * time.Minute
	var tpl *annon.CatalogTemplate

	if startTime == nil {
		log.WithField("template", catalogID).Warn("Received generation request without startTime")
		return e.New("Start time cannot be empty")
	}

	clog = clog.WithField("template", catalogID)
	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"startTime": startTime, "endTime": endTime, "dayInterval": dayInterval, "nightInterval": nightInterval, "params": params}).Debug("Generating announcement from template.")
	}

	if tpl, err = (*s.cat).Get(catalogID); err != nil {
		clog.WithError(err).Error("Could not retrieve catalog template.")
		return err
	}
	if tpl == nil {
		clog.Debug("Template not found")
		return errors.NewError("Template not found.", errors.NotFound)
	}
	clog.Debug("Loaded requested template.")

	now := time.Now()
	timezonedStart := (*startTime).In(config.Timezone)
	var timezonedEnd time.Time
	if endTime != nil {
		timezonedEnd = (*endTime).In(config.Timezone)
	}
	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"startTime": timezonedStart, "endTime": timezonedEnd}).Debug("Aligned time parameters with the current config.Timezone.")
	}
	var annons *annon.Announcement
	if annons, err = (*s.gen).FromTemplate(tpl, &timezonedStart, &timezonedEnd, &dayInt, &nightInt, &now, params, 10); err != nil {
		clog.WithError(err).Error("Could not generate announcements for template.")
		return err
	}
	if annons == nil {
		clog.Info("Announcement was not created. Ignoring.")
		return nil
	}
	// TODO refactor introducing zoneID properly
	st, _ := strconv.Atoi(zoneID)
	annons.StationID = st
	var id string
	var duration int
	var approxLen bool

	clog.Debug("Generating announcement audio.")

	if id, _, duration, approxLen, err = (*s.audio).Generate(annons.Text.TtsText); err != nil {
		clog.WithError(err).Error("Could not generate audio file.")
		return err
	}
	aud := annon.Audio{Duration: duration, FileID: id, ApproxLen: approxLen}
	annons.Audio = &aud

	if log.GetLevel() >= log.DebugLevel {
		clog.WithField("announcement", annons).Debug("Saving generated announcement.")
	}

	if _, err = (*s.store).Save(annons); err != nil {
		clog.WithError(err).Error("Error saving announcements for template.")
		return err
	}
	return nil
}

func (s *announcement) GetCatalog() (catalog map[string][]*annon.CatalogTemplate, err error) {
	log.WithFields(log.Fields{"logger": "annon.service", "method": "GetCatalog"}).
		Debug("Loading catalog contents.")
	return (*s.cat).GetAll()
}

func (s *announcement) GetTemplate(templateID string) (*annon.CatalogTemplate, error) {
	if log.GetLevel() > log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annon.service", "method": "GetTemplate", "template": templateID}).
			Debug("Loading catalog template.")
	}
	return (*s.cat).Get(templateID)
}

func (s *announcement) DeleteTemplate(templateID string) error {
	if log.GetLevel() > log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annon.service", "method": "DeleteTemplate", "template": templateID}).
			Debug("Deleting catalog template.")
	}
	return (*s.cat).Delete(templateID)
}

func (s *announcement) SaveTemplate(templ *annon.CatalogTemplate) (*annon.CatalogTemplate, error) {
	if log.GetLevel() > log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annon.service", "method": "SaveTemplate", "template": fmt.Sprintf("%+v", *templ)}).
			Debug("Deleting catalog template.")
	}
	return (*s.cat).Save(templ)
}
