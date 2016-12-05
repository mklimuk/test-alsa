package annon

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/mklimuk/husar/l10n"
	"github.com/mklimuk/husar/train"

	"golang.org/x/text/language"
)

func init() {
	dep := NewDeparture()
	producers[dep.Name()] = &dep
}

/*NewDeparture is a constructor of the departure producer*/
func NewDeparture() Producer {
	a := new(departure)
	a.name = Departure
	return Producer(a)
}

// departure implements TemplateHandler interface for departure announcements
type departure struct {
	producer
}

// Required implements TemplateHandler Required method
func (d *departure) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	if t == nil && old != nil {
		isRequired = false
		// to delete
		needsUpdate = true
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.departure", "method": "Required", "train": old.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Removing train.")
		}
		return
	}
	if t.Route == nil {
		isRequired = false
		needsUpdate = false
		log.WithFields(log.Fields{"logger": "annongen.producer.departure", "method": "Required", "train": t.ID}).
			Warn("Route not present in train struct.")
		return
	}
	// if there is a subroute this is a replacement bus so no departure
	if t.Route.CurrentStationOnSubroute {
		isRequired = false
		needsUpdate = false
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.departure", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Station is on a subroute. Train departure announcement should be replaced by bus arrival announcement.")
		}
		return
	}
	isRequired = t.Departure != nil && t.StationID != t.Route.EndStation.ID
	if !isRequired {
		needsUpdate = false
	} else {
		needsUpdate = d.needsUpdate(t, old)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
			Debug("Departure announcement required check results")
	}
	return
}

func (*departure) needsUpdate(current *train.Train, old *train.Train) bool {
	if old == nil {
		return true
	}
	return train.GetDelay(current) != train.GetDelay(old) ||
		!train.CompareEvents(train.GetDeparture(current), train.GetDeparture(old)) ||
		!train.CompareI18n(current, old, train.DepartureText) ||
		!train.CompareOverrides(current, old)
}

/*BuildParams builds template handler parameters for arrival handler*/
func (d *departure) BuildParams(t *train.Train, lang language.Tag) (*TemplateParams, *TemplateParams, bool, error) {

	dep := train.GetDeparture(t)

	if dep == nil {
		err := fmt.Errorf("Cannot build departure announcement parameters. Departure event not found for train %v", t)
		return nil, nil, false, err
	}
	_, through := train.FollowingStations(t.Route, true)
	first := (t.Route.StartStation.ID == t.StationID)
	params := TemplateParams{
		"Category":  l10n.FromMetaDictionary(train.GetCategory(t), l10n.Genitive, lang),
		"Carrier":   l10n.FromMetaDictionary(t.Carrier, l10n.Locative, lang),
		"Name":      t.Name,
		"NameLower": strings.ToLower(t.Name),
		"From":      t.Route.StartStation.Name,
		"To":        t.Route.EndStation.Name,
		"By":        through,
		"Track":     dep.Track,
		"TrackTxt":  strconv.Itoa(dep.Track),
		"Platform":  dep.Platform,
		"First":     first,
		"Delayed":   train.GetDelay(t) > 0,
	}
	if res, val := train.HasI18nOption(t, train.DepartureText, lang); res == true {
		params["Custom"] = val
	}
	ttsParams := d.MapParams(params, lang)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.departure", "method": "BuildParams", "params": fmt.Sprintf("%+v", params), "ttsParams": fmt.Sprintf("%+v", ttsParams)}).
			Debug("Parameters for speech generator")
	}
	return &params, &ttsParams, true, nil
}

func (*departure) GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time) {
	scheduled := t.Departure.Time
	if t.Carrier == "IC" {
		scheduled = scheduled.Add(time.Duration(train.GetDelay(t)*60-30) * time.Second)
	} else {
		scheduled = scheduled.Add(time.Duration(train.GetDelay(t)-2) * time.Minute)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.departure", "method": "GetTime", "expected": scheduled}).
			Debug("Announcement expected time.")
	}
	return append(events, scheduled), scheduled, scheduled
}

func (*departure) MapParams(params TemplateParams, lang language.Tag) TemplateParams {
	mapped := CopyParams(params)
	mapped["Platform"] = l10n.NumToText(l10n.RomanToInt(params["Platform"].(string)), l10n.Locative, lang)
	mapped["TrackTxt"] = l10n.NumToText(params["Track"].(int), l10n.Adjective, lang)
	return mapped
}
