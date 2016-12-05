package annon

import "time"

/*
Category announcements category
*/
type Category string

/*
Available categories
*/
const (
	Train   Category = "train"
	Special Category = "special"
)

/*
Available announcement types
*/
const (
	Arrival       string = "arrival"
	BusArrival    string = "busArrival"
	Departure     string = "departure"
	BusDeparture  string = "busDeparture"
	Delay         string = "delay"
	Setup         string = "setup"
	BusSetup      string = "busSetup"
	Stop          string = "stop"
	Info          string = "info"
	InfoWithPhone string = "info-phone"
)

//Priority is a playback priority of an announcement
type Priority int

// Predefined priorities
const (
	P1 Priority = iota
	P2
	P3
	P4
	P5
)

// Announcement represents a voice announcement broadcasted by the system
type Announcement struct {
	ID        string      `gorethink:"id,omitempty" json:"id"`
	Time      []time.Time `gorethink:"time" json:"time"`
	First     time.Time   `gorethink:"first" json:"first"`
	Last      time.Time   `gorethink:"last" json:"last"`
	Category  Category    `gorethink:"category" json:"category"`
	Lang      string      `gorethink:"lang" json:"lang"`
	Type      string      `gorethink:"type" json:"type"`
	Priority  Priority    `gorethink:"priority" json:"priority"`
	TrainID   string      `gorethink:"trainId" json:"trainId,omitempty"`
	StationID int         `gorethink:"stationId" json:"stationId"`
	Text      *Text       `gorethink:"text" json:"text"`
	Audio     *Audio      `gorethink:"audio" json:"audio"`
	Autoplay  bool        `gorethink:"autoplay" json:"autoplay"`
}

// Text properties of an announcement
type Text struct {
	HumanText string `gorethink:"humanText" json:"humanText"`
	HTMLText  string `gorethink:"htmlText" json:"htmlText"`
	TtsText   string `gorethink:"ttsText" json:"ttsText"`
}

// Audio properties of an announcement, duration is in seconds
type Audio struct {
	FileID    string `gorethink:"fileId" json:"fileId"`
	Duration  int    `gorethink:"duration" json:"duration"`
	ApproxLen bool   `gorethink:"approxLen" json:"approxLen"`
}

// Settings can be manipulated by the user
type Settings struct {
	IsBlocked bool `gorethink:"isBlocked" json:"isBlocked"`
}

// Changes represents changefeeds data from database
type Changes struct {
	NewVal *Announcement `gorethink:"new_val,omitempty"`
	OldVal *Announcement `gorethink:"old_val,omitempty"`
}
