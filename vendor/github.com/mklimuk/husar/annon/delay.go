package annon

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/mklimuk/husar/l10n"
	"github.com/mklimuk/husar/train"

	"golang.org/x/text/language"
)

func init() {
	de := NewDelay()
	producers[de.Name()] = &de
}

/*NewDelay is the delay announcement producer's constructor*/
func NewDelay() Producer {
	a := new(delay)
	a.name = Delay
	return Producer(a)
}

// delayed implements TemplateHandler interface for delayed announcements
type delay struct {
	producer
}

// Required implements TemplateHandler Required method
func (a *delay) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	if t == nil && old != nil {
		isRequired = false
		// to delete
		needsUpdate = true
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "Required", "train": old.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Removing train.")
		}
		return
	}
	if t.Route == nil {
		isRequired = false
		needsUpdate = false
		log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "Required", "train": t.ID}).
			Warn("Route not present in train struct.")
		return
	}
	event := train.GetReferenceEvent(t)
	if event == nil {
		err = fmt.Errorf("Reference event not found for train")
		log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "Required", "train": fmt.Sprintf("%+v", *t)}).
			WithError(err).Error("Could not check if producer is required")
		return false, false, err
	}

	isRequired = train.GetDelay(t) > 0
	needsUpdate = train.GetDelay(t) != train.GetDelay(old) ||
		!train.CompareI18n(t, old, train.DelayText) ||
		!train.CompareOverrides(t, old)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
			Debug("Delay announcement required check results")
	}
	return
}

// buildParams builds template handler parameters for arrival handler
func (a *delay) BuildParams(t *train.Train, lang language.Tag) (*TemplateParams, *TemplateParams, bool, error) {
	_, through := train.FollowingStations(t.Route, true)
	last := (t.Route.EndStation.ID == t.StationID)
	first := (t.Route.StartStation.ID == t.StationID)
	event := train.GetReferenceEvent(t)

	params := TemplateParams{
		"Category":  l10n.FromMetaDictionary(train.GetCategory(t), l10n.Genitive, lang),
		"Carrier":   l10n.FromMetaDictionary(t.Carrier, l10n.Locative, lang),
		"Name":      t.Name,
		"From":      t.Route.StartStation.Name,
		"To":        t.Route.EndStation.Name,
		"By":        through,
		"Last":      last,
		"First":     first,
		"Delay":     train.GetDelay(t),
		"Scheduled": event.Time.Format(TimeFormat(lang)),
	}
	if res, val := train.HasI18nOption(t, train.DelayText, lang); res == true {
		params["Custom"] = val
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "BuildParams", "params": fmt.Sprintf("%+v", params)}).
			Debug("Parameters for speech generator")
	}
	return &params, &params, true, nil
}

func (a *delay) GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time) {
	scheduled := train.GetReferenceScheduledTime(t)
	diff := scheduled.Sub(*now)
	// if the train should have leaved we add immediate announcement and then one every thirty minutes
	if diff < 0 {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.delay", "method": "GetTime", "train": t.ID, "scheduled": scheduled, "now": now}).
				Debug("We are already after the train's scheduled departure. Adding events.")
		}
		// announce the delay in a few seconds
		events = append(events, now.Add(30*time.Second))
		// we repeat the announcement until 5 minutes before arrival
		events = a.appendUntilEstimated(events, now.Add(30*time.Minute), scheduled.Add(time.Duration(train.GetDelay(t)-5)*time.Minute))
	} else {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "producer.delay", "train": t.ID, "scheduled": scheduled, "now": now}).
				Debug("We are before train's scheduled departure. Adding events.")
		}
		events = append(events, scheduled.Add(-10*time.Minute), scheduled.Add(-5*time.Minute))
		// if the delay is bigger than 5 minutes we also append an announcement at the expected arrival/departure time
		if t.Delay > 5 {
			events = a.appendUntilEstimated(events, scheduled, scheduled.Add(time.Duration(train.GetDelay(t)-5)*time.Minute))
		}
	}
	return events, events[0], events[len(events)-1]
}

func (a *delay) appendUntilEstimated(events []time.Time, start time.Time, end time.Time) []time.Time {
	current := start
	for current.Before(end) {
		events = append(events, current)
		current = current.Add(30 * time.Minute)
	}
	return events
}
