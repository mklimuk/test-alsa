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
	arr := NewArrival()
	producers[arr.Name()] = &arr
}

/*NewBusArrival is a constructor of the arrival producer*/
func NewBusArrival() Producer {
	a := new(busArrival)
	a.name = BusArrival
	return Producer(a)
}

// busArrival implements TemplateHandler interface for bus arrival announcements
type busArrival struct {
	producer
}

// Required implements TemplateHandler Required method
func (a *busArrival) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	if t == nil && old != nil {
		isRequired = false
		// to delete
		needsUpdate = true
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.busArrival", "method": "Required", "train": old.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Removing train.")
		}
		return
	}
	if t.Route == nil {
		isRequired = false
		needsUpdate = false
		log.WithFields(log.Fields{"logger": "annongen.producer.busArrival", "method": "Required", "train": t.ID}).
			Warn("Route not present in train struct.")
		return
	}
	isRequired = t.Arrival != nil && t.Route.CurrentStationOnSubroute
	if !isRequired {
		needsUpdate = false
	} else {
		needsUpdate = a.needsUpdate(t, old)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busArrival", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
			Debug("Bus arrival announcement required check results")
	}
	return
}

func (*busArrival) needsUpdate(current *train.Train, old *train.Train) bool {
	if old == nil {
		return true
	}
	return train.GetDelay(current) != train.GetDelay(old) ||
		!train.CompareI18n(current, old, train.BusArrivalText) ||
		!train.CompareOverrides(current, old)
}

// buildParams builds template handler parameters for arrival handler
func (a *busArrival) BuildParams(t *train.Train, lang language.Tag) (*TemplateParams, *TemplateParams, bool, error) {

	params := TemplateParams{
		"Category":      l10n.FromMetaDictionary(train.GetCategory(t), l10n.Genitive, lang),
		"Carrier":       l10n.FromMetaDictionary(t.Carrier, l10n.Locative, lang),
		"SubrouteStart": t.Route.SubrouteStart.Name,
		"Name":          t.Name,
		"From":          t.Route.StartStation.Name,
		"To":            t.Route.EndStation.Name,
		"Delayed":       train.GetDelay(t) > 0,
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busArrival", "method": "BuildParams", "params": fmt.Sprintf("%+v", params)}).
			Debug("Parameters for speech generator")
	}
	return &params, &params, false, nil
}

func (*busArrival) GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time) {
	expected := t.Arrival.Time.Add(time.Duration(train.GetDelay(t)-5) * time.Minute)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.busArrival", "method": "GetTime", "expected": expected}).
			Debug("Announcement expected time.")
	}
	return append(events, expected), expected, expected
}
