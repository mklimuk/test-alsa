package msg

import (
	"encoding/json"
	"time"

	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

//Type tells us what to modify
type Type string

//available event types
const (
	Arrival   Type = "arrival"
	Departure Type = "departure"
	All       Type = "all"
)

//TimetableRequest is the timetable information request struct
type TimetableRequest struct {
	StationID int       `json:"stationId"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

//TimetableResponse is the timetable information response struct
type TimetableResponse struct {
	StationID int               `json:"stationId"`
	Settings  *[]train.Settings `json:"settings"`
	Timetable []*train.Train    `json:"timetable"`
}

//TrainUpdate is the train update request message struct
type TrainUpdate struct {
	RequestMessage
	TrainID   string            `json:"trainId"`
	StationID string            `json:"stationId"`
	Mode      *train.UpdateMode `json:"mode,omitempty"`
	Event     *Event            `json:"event,omitempty"`
}

//AudioSettingsUpdate is the audio settings update message struct
type AudioSettingsUpdate struct {
	TrainID   string                            `json:"trainId"`
	Settings  map[train.AnnonOption]string      `json:"settings"`
	Overrides map[string]string                 `json:"overrides"`
	Lang      map[string]train.LanguageSettings `json:"lang"`
}

//Event represents simplified timetable events in messages
type Event struct {
	Track    *int    `json:"track,omitempty"`
	Platform *string `json:"platform,omitempty"`
	Delay    *int    `json:"delay,omitempty"`
	Type     Type    `json:"eventType,omitempty"`
}

//ParseTrainUpdate parses json train update message
func ParseTrainUpdate(req string) (*TrainUpdate, error) {
	msg := new(TrainUpdate)
	var err error
	if err = json.Unmarshal([]byte(req), msg); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "parseTrainUpdate"}).
			WithError(err).Error("Error unmarshalling request")
		return nil, err
	}
	return msg, nil
}

//ParseSettingsUpdate parses json settings update message
func ParseSettingsUpdate(req string) (*AudioSettingsUpdate, error) {
	msg := new(AudioSettingsUpdate)
	var err error
	if err = json.Unmarshal([]byte(req), msg); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.timetable", "method": "parseSettingsUpdate"}).
			WithError(err).Error("Error unmarshalling request")
		return nil, err
	}
	return msg, nil
}
