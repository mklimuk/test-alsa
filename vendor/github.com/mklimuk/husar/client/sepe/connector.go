package sepe

import "net/url"

//Connector handles communication with SEPE service
type Connector interface {
	CheckHealth() (bool, error)
	ProcessScheduledTimetables() error
	ProcessLiveTimetables() error
}

//NewConnector is SEPE connector constructor
func NewConnector(baseURL string, monitorURL string) Connector {
	var err error
	var base, monitor *url.URL
	if base, err = url.Parse(baseURL); err != nil {
		panic(err)
	}
	if monitor, err = url.Parse(monitorURL); err != nil {
		panic(err)
	}
	c := connector{baseURL: base, monitorBaseURL: monitor}
	return Connector(&c)
}

type connector struct {
	baseURL        *url.URL
	monitorBaseURL *url.URL
}
