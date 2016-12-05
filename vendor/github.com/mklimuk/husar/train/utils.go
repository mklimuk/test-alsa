package train

import (
	"time"

	"golang.org/x/text/language"

	log "github.com/Sirupsen/logrus"
)

// HasService verifies if a service is present for a given train and returns
// the list of carriages that propose the service
func HasService(train *Train, st ServiceType) (present bool, carriage *[]string) {
	if train.Services == nil {
		return false, nil
	}
	for _, service := range *train.Services {
		if service.ID == string(st) {
			return true, &service.Carriage
		}
	}
	return false, nil
}

// HasAnnonOption checks if an announcement option is set on a train and returns its value
func HasAnnonOption(train *Train, opt AnnonOption) (present bool, value string) {
	if train == nil {
		log.Error("Received a nil train. This should never happen.")
		return false, ""
	}
	if train.Settings == nil || train.Settings.Audio == nil {
		return false, ""
	}
	for key, value := range *train.Settings.Audio {
		if key == opt {
			return true, value
		}
	}
	return false, ""
}

// HasI18nOption checks if an announcement option is set on a train and returns its value
func HasI18nOption(train *Train, opt AnnonOption, lang language.Tag) (present bool, value string) {
	if train == nil {
		log.Error("Received a nil train. This should never happen.")
		return false, ""
	}
	if train.Settings == nil || train.Settings.Lang == nil {
		return false, ""
	}
	options := (*train.Settings.Lang)[lang.String()].I18n
	if options == nil {
		return false, ""
	}
	value, present = (*options)[opt]
	return
}

// HasI18n checks if i18n is enabled for a given train
func HasI18n(train *Train) (bool, *map[string]LanguageSettings) {
	if train == nil {
		log.Error("Received a nil train. This should never happen.")
		return false, nil
	}
	if train.Settings == nil || train.Settings.Lang == nil {
		return false, nil
	}
	return true, train.Settings.Lang
}

//CompareAudioSettings returns true if a given audio setting is equal for two trains.
func CompareAudioSettings(current *Train, old *Train, prop AnnonOption) bool {
	if old == nil {
		return current == nil
	}
	if current == nil {
		return false
	}
	if current.Settings == nil {
		return old.Settings == nil
	}
	if old.Settings == nil {
		return false
	}
	if current.Settings.Audio == nil {
		return old.Settings.Audio == nil
	}
	if old.Settings.Audio == nil {
		return false
	}
	return (*current.Settings.Audio)[prop] == (*old.Settings.Audio)[prop]
}

//CompareI18n returns true if i18n audio setting is equal for two trains.
func CompareI18n(current *Train, old *Train, prop AnnonOption) bool {
	if old == nil {
		return current == nil
	}
	if current == nil {
		return false
	}

	if current.Settings == nil {
		return old.Settings == nil
	}
	if old.Settings == nil {
		return false
	}

	if current.Settings.Lang == nil {
		return old.Settings.Lang == nil
	}
	if old.Settings.Lang == nil {
		return false
	}
	for k, v := range *current.Settings.Lang {
		other := (*old.Settings.Lang)[k]
		if v.Enabled != other.Enabled {
			return false
		}
		if (v.I18n == nil) != (other.I18n == nil) {
			return false
		}
		if v.I18n != nil {
			for key, txt := range *v.I18n {
				if txt != (*other.I18n)[key] {
					return false
				}
			}
		}
	}
	return true
}

//CompareOverrides returns true if overriden train settings are equal for two trains.
func CompareOverrides(current *Train, old *Train) bool {
	if old == nil {
		return current == nil
	}
	if current == nil {
		return false
	}
	if current.Settings == nil {
		return old.Settings == nil
	}
	if old.Settings == nil {
		return false
	}
	if current.Settings.Overrides == nil {
		return old.Settings.Overrides == nil
	}
	if old.Settings.Overrides == nil {
		return false
	}
	for k, v := range *current.Settings.Overrides {
		if v != (*old.Settings.Overrides)[k] {
			return false
		}
	}
	return true
}

// FollowingStations returns a slice of all stations that are located after the current station on route.
// It also returns a concatenated string of those station's names.
func FollowingStations(route *Route, stripLast bool) (stations []Station, through string) {
	first := true
	for _, station := range route.Stations {
		if station.NumberOnRoute > route.CurrentStationPositionOnRoute &&
			station.Important &&
			station.ID != route.StartStation.ID &&
			(station.ID != route.EndStation.ID || !stripLast) {
			stations = append(stations, station)
			if first == true {
				first = false
			} else {
				through += ", "
			}
			through += station.Name
		}
	}
	return stations, through
}

//BuildSettingsMap produces a map from a slice of Settings
func BuildSettingsMap(set *[]Settings) map[string]Settings {
	if set == nil {
		return map[string]Settings{}
	}
	settings := make(map[string]Settings, len(*set))
	for _, s := range *set {
		settings[s.ID] = s
	}
	return settings
}

//Copy one train struct into another
func Copy(source *Train, c *Train) *Train {
	*c = *source
	if source.Arrival != nil {
		a := new(TimetableEvent)
		*a = *(source.Arrival)
		c.Arrival = a
	}
	if source.Departure != nil {
		d := new(TimetableEvent)
		*d = *(source.Departure)
		c.Departure = d
	}
	if source.Route != nil {
		r := new(Route)
		*r = *(source.Route)
		c.Route = r
		s := make([]Station, len(source.Route.Stations))
		copy(s, r.Stations)
		r.Stations = s
		c.Route = r
	}
	return c
}

// GetArrival returns arrival event for a train taking into account the fact that it may have manual settings activated
func GetArrival(t *Train) *TimetableEvent {
	if t == nil {
		return nil
	}
	if t.Settings != nil && t.Settings.Arrival != nil {
		return t.Settings.Arrival
	}
	return t.Arrival
}

// GetDeparture returns departure event for a train taking into account the fact that it may have manual settings activated
func GetDeparture(t *Train) *TimetableEvent {
	if t == nil {
		return nil
	}
	if t.Settings != nil && t.Settings.Departure != nil {
		return t.Settings.Departure
	}
	return t.Departure
}

// GetReferenceEvent returns reference event for a train taking into account the fact that it may have manual settings activated
func GetReferenceEvent(t *Train) *TimetableEvent {
	var res *TimetableEvent
	if res = GetArrival(t); res != nil {
		return res
	}
	return GetDeparture(t)
}

// GetDelay returns delay for a train taking into account the fact that it may have manual settings activated
func GetDelay(t *Train) int {
	if t == nil {
		return 0
	}
	if t.Settings != nil && t.Settings.Mode == Manual {
		return t.Settings.Delay
	}
	return t.Delay
}

// GetFirstLiveEvent calculates first live event based on train settings
func GetFirstLiveEvent(t *Train) *time.Time {
	delay := GetDelay(t)
	var event time.Time
	if t.Arrival != nil {
		event = t.Arrival.Time.Add(time.Duration(delay) * time.Minute)
	}
	if t.Departure != nil {
		event = t.Departure.Time.Add(time.Duration(delay) * time.Minute)
	}
	return &event
}

// CompareEvents returns true if two events are equal
func CompareEvents(a *TimetableEvent, b *TimetableEvent) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.Platform == b.Platform && a.Track == b.Track && a.Time.Equal(b.Time)
}

//GetReferenceScheduledTime returns scheduled time of a reference event (ignoring manual settings)
func GetReferenceScheduledTime(t *Train) time.Time {
	if t.Arrival != nil {
		return t.Arrival.Time
	}
	return t.Departure.Time
}

//ModeChanged returns true if update mode is different between two Train instances
func ModeChanged(a *Train, b *Train) bool {
	if a.Settings == nil {
		if b.Settings != nil {
			return true
		}
		return false
	}
	if b.Settings == nil {
		return true
	}
	return a.Settings.Mode == b.Settings.Mode
}

// GetCategory returns train's category taking into account potential overrides
func GetCategory(t *Train) string {
	if t == nil {
		return ""
	}
	if t.Settings != nil && t.Settings.Overrides != nil && (*t.Settings.Overrides)[Category] != "" {
		return (*t.Settings.Overrides)[Category]
	}
	return t.Category
}
