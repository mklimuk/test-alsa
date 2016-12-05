package queue

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type queueContent struct {
	QueueID   string    `json:"queueId"`
	Mute      bool      `json:"mute"`
	Events    *[]*Event `json:"events"`
	NextEvent *Event    `json:"next"`
}

type queueMessage struct {
	QueueID string `json:"queueId"`
	EventID string `json:"eventId"`
}

func parseQueueMessage(req string) *queueMessage {
	msg := new(queueMessage)
	var err error
	if err = json.Unmarshal([]byte(req), msg); err != nil {
		log.WithFields(log.Fields{"logger": "lcs.queue", "method": "parseQueueMessage"}).
			WithError(err).Error("Error unmarshalling request")
		return nil
	}
	return msg
}
