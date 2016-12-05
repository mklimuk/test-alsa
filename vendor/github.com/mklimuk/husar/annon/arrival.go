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
	arr := NewArrival()
	producers[arr.Name()] = &arr
}

/*NewArrival is a constructor of the arrival producer*/
func NewArrival() Producer {
	a := new(arrival)
	a.name = Arrival
	return Producer(a)
}

// arrival implements TemplateHandler interface for arrival announcements
type arrival struct {
	producer
}

// Required implements annon.TemplateHandler Required method
func (a *arrival) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	if t == nil && old != nil {
		isRequired = false
		// to delete
		needsUpdate = true
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "Required", "train": old.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Removing train.")
		}
		return
	}
	if t.Route == nil {
		isRequired = false
		needsUpdate = false
		log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "Required", "train": t.ID}).
			Warn("Route not present in train struct.")
		return
	}
	// if there is a subroute this is a replacement bus so no departure
	if t.Route.CurrentStationOnSubroute {
		isRequired = false
		needsUpdate = false
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
				Debug("Station is on a subroute. Train arrival announcement should be replaced by bus arrival announcement.")
		}
		return
	}
	isRequired = t.Arrival != nil && t.StationID != t.Route.StartStation.ID
	if !isRequired {
		needsUpdate = false
	} else {
		needsUpdate = a.needsUpdate(t, old)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "Required", "train": t.ID, "required": isRequired, "needsUpdate": needsUpdate}).
			Debug("Arrival announcement required check results")
	}
	return
}

func (*arrival) needsUpdate(current *train.Train, old *train.Train) bool {
	if old == nil {
		return true
	}
	return train.GetDelay(current) != train.GetDelay(old) ||
		!train.CompareEvents(train.GetArrival(current), train.GetArrival(old)) ||
		!train.CompareI18n(current, old, train.ArrivalText) ||
		!train.CompareOverrides(current, old)
}

// buildParams builds template handler parameters for arrival handler
func (a *arrival) BuildParams(t *train.Train, lang language.Tag) (*TemplateParams, *TemplateParams, bool, error) {

	// we get the current arrival
	arr := train.GetArrival(t)

	if arr == nil {
		err := fmt.Errorf("Cannot build arrival announcement parameters. Arrival event not found for train %v", t)
		return nil, nil, false, err
	}

	_, through := train.FollowingStations(t.Route, true)
	last := (t.Route.EndStation.ID == t.StationID)
	params := TemplateParams{
		"Category":  l10n.FromMetaDictionary(train.GetCategory(t), l10n.Genitive, lang),
		"Carrier":   l10n.FromMetaDictionary(t.Carrier, l10n.Locative, lang),
		"Name":      t.Name,
		"NameLower": strings.ToLower(t.Name),
		"From":      t.Route.StartStation.Name,
		"To":        t.Route.EndStation.Name,
		"By":        through,
		"Track":     arr.Track,
		"TrackTxt":  strconv.Itoa(arr.Track),
		"Platform":  arr.Platform,
		"Last":      last,
		"Delayed":   train.GetDelay(t) > 0,
	}
	if last == false {
		params["Departure"] = t.Departure.Time.Format(TimeFormat(lang))
	}
	if res, _ := train.HasService(t, train.REZO); res == true {
		params["Resa"] = true
	}
	if res, car := train.HasService(t, train.KOND); res == true {
		params["Kond"] = (*car)[0]
	}
	if res, val := train.HasI18nOption(t, train.ArrivalText, lang); res == true {
		params["Custom"] = val
	}
	ttsParams := a.MapParams(params, lang)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "BuildParams", "params": fmt.Sprintf("%+v", params), "ttsParams": fmt.Sprintf("%+v", ttsParams)}).
			Debug("Parameters for speech generator")
	}
	return &params, &ttsParams, true, nil
}

func (*arrival) GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time) {
	expected := t.Arrival.Time.Add(time.Duration(train.GetDelay(t)-5) * time.Minute)
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.producer.arrival", "method": "GetTime", "expected": expected}).
			Debug("Announcement expected time.")
	}
	return append(events, expected), expected, expected
}

func (*arrival) MapParams(params TemplateParams, lang language.Tag) TemplateParams {
	mapped := CopyParams(params)
	mapped["Platform"] = l10n.NumToText(l10n.RomanToInt(params["Platform"].(string)), l10n.Locative, lang)
	mapped["TrackTxt"] = l10n.NumToText(params["Track"].(int), l10n.Genitive, lang)
	return mapped
}
