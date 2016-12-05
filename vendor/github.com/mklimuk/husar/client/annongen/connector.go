package annongen

import (
	"net/url"
	"time"
)

//Connector is annongen control API
type Connector interface {
	CheckHealth() (bool, error)
	ProcessSchedule(from *time.Time, to *time.Time) error
	EnableTimetableSync() error
}

//NewConnector is a connector constructor
func NewConnector(baseURL string) Connector {
	var err error
	var u *url.URL
	if u, err = url.Parse(baseURL); err != nil {
		panic(err)
	}
	c := connector{baseURL: u}
	return Connector(&c)
}

type connector struct {
	baseURL *url.URL
}
