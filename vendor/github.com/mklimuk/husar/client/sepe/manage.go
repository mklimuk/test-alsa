package sepe

import (
	"errors"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

const (
	scheduledPath string = "/timetable/scheduled/update"
	livePath      string = "/timetable/live/update"
)

func (c *connector) ProcessScheduledTimetables() error {
	u := new(url.URL)
	*u = *c.baseURL
	u.Path += scheduledPath
	scheduledURL := u.String()
	clog := log.WithFields(log.Fields{"logger": "init.sepe", "method": "ProcessScheduledTimetables"})
	clog.WithField("url", scheduledURL).Info("Requesting scheduled timetable processing.")
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(scheduledURL); err != nil {
		clog.WithError(err).WithField("url", scheduledURL).Info("Error while calling SEPE scheduled timetable processing URL.")
		return err
	}
	st := resp.StatusCode
	if st != http.StatusOK {
		clog.WithField("status", st).Info("Process scheduled timetable response status is not OK.")
		return errors.New("Process scheduled timetable response status is not OK")
	}
	return nil
}

func (c *connector) ProcessLiveTimetables() error {
	u := new(url.URL)
	*u = *c.baseURL
	u.Path += livePath
	liveURL := u.String()
	clog := log.WithFields(log.Fields{"logger": "init.sepe", "method": "ProcessLiveTimetables"})
	clog.WithField("url", liveURL).Info("Requesting live timetable processing.")
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(liveURL); err != nil {
		clog.WithField("url", liveURL).WithError(err).Info("Error while calling SEPE live timetable processing URL.")
		return err
	}
	st := resp.StatusCode
	if st != http.StatusOK {
		clog.WithField("status", st).Info("Process live timetable response status is not OK.")
		return errors.New("Process live timetable response status is not OK")
	}
	return nil
}
