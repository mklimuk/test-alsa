package tts

import "time"

//GenerateRequest is a struct used as TTS requests payload
type GenerateRequest struct {
	ID   string    `json:"id"`
	Time time.Time `json:"time"`
	Text string    `json:"text"`
}
