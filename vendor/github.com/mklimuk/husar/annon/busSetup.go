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
	dep := NewBusSetup()
	producers[dep.Name()] = &dep
}

/*NewBusSetup is a constructor of the setup producer*/
func NewBusSetup() Producer {
	a := new(busSetup)
	a.name = BusSetup
	return Producer(a)
}

// busSetup implements TemplateHandler interface for setup announcements
type busSetup struct {
	producer
}

// Required implements TemplateHandler Required method
func (s *busSetup) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	if t == nil && old != nil {
		isRequired = false
		// to delete
		needsUpdate = true
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.busSetup", "method": "Required", "train": old.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Removing train.")
		}
		return
	}
	if t.Route == nil {
		isRequired = false
		needsUpdate = false
		log.WithFields(log.Fields{"logger": "annongen.producer.busSetup", "method": "Required", "train": t.ID}).
			Warn("Route not present in train struct.")
		return
	}
	isRequired = t.Departure != nil && t.Route.CurrentStationOnSubroute
	if !isRequired {
		needsUpdate = false
	} else {
		needsUpdate = s.needsUpdate(t, old)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busSetup", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
			Debug("Bus arrival announcement required check results")
	}
	return
}

func (*busSetup) needsUpdate(current *train.Train, old *train.Train) bool {
	if old == nil {
		return true
	}
	return train.GetDelay(current) != train.GetDelay(old) ||
		!train.CompareI18n(current, old, train.BusSetupText) ||
		!train.CompareOverrides(current, old)
}

// buildParams builds template handler parameters for arrival handler
func (s *busSetup) BuildParams(t *train.Train, lang language.Tag) (*TemplateParams, *TemplateParams, bool, error) {
	dep := train.GetDeparture(t)
	if dep == nil {
		err := fmt.Errorf("Cannot build setup announcement parameters. Setup event not found for train %v", t)
		return nil, nil, false, err
	}
	params := TemplateParams{
		"Category":    l10n.FromMetaDictionary(train.GetCategory(t), l10n.Genitive, lang),
		"Carrier":     l10n.FromMetaDictionary(t.Carrier, l10n.Locative, lang),
		"Name":        t.Name,
		"From":        t.Route.StartStation.Name,
		"To":          t.Route.EndStation.Name,
		"SubrouteEnd": t.Route.SubrouteEnd.Name,
		"Departure":   t.Departure.Time.Format(TimeFormat(lang)),
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busSetup", "method": "BuildParams", "params": fmt.Sprintf("%+v", params)}).
			Debug("Parameters for speech generator")
	}
	return &params, &params, false, nil
}

func (*busSetup) GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time) {
	scheduled := t.Departure.Time.Add(time.Duration(train.GetDelay(t)-150) * time.Second)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busSetup", "method": "GetTime", "expected": scheduled}).
			Debug("Announcement expected time.")
	}
	return append(events, scheduled), scheduled, scheduled
}
