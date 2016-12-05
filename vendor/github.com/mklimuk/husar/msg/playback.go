package msg

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

//VolMessage is the volume modification request struct
type VolMessage struct {
	QueueID string `json:"queueId"`
	Volume  int    `json:"volume"`
}

//TriggerMessage is the announcement trigger request struct
type TriggerMessage struct {
	QueueID string `json:"queueId"`
	EventID string `json:"eventId"`
}

//ParseVolumeMessage parses VolMessage from json string
func ParseVolumeMessage(req string) *VolMessage {
	msg := new(VolMessage)
	var err error
	if err = json.Unmarshal([]byte(req), msg); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.playback", "method": "parseVolumeMessage"}).
			WithError(err).Error("Error unmarshalling request")
		return nil
	}
	return msg
}

//ParseTriggerMessage parses TriggerMessage from json string
func ParseTriggerMessage(req string) *TriggerMessage {
	msg := new(TriggerMessage)
	var err error
	if err = json.Unmarshal([]byte(req), msg); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.playback", "method": "parseTriggerMessage"}).
			WithError(err).Error("Error unmarshalling request")
		return nil
	}
	return msg
}
