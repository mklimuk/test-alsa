package annongen

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	syncPath    string = "/timetable/sync/enable"
	processPath string = "/timetable/process"
)

func (c *connector) ProcessSchedule(from *time.Time, to *time.Time) error {
	u := new(url.URL)
	*u = *c.baseURL
	u.Path += processPath
	params := url.Values{}
	params.Add("from", from.Format(time.RFC3339))
	params.Add("to", to.Format(time.RFC3339))
	u.RawQuery = params.Encode()
	processURL := u.String()
	clog := log.WithFields(log.Fields{"logger": "init.annongen", "method": "ProcessSchedule", "url": u.RequestURI()})
	clog.WithField("url", processURL).Info("Requesting scheduled timetable processing.")
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(processURL); err != nil {
		clog.WithField("url", processURL).WithError(err).Info("Error while calling annongen processing URL.")
		return err
	}
	st := resp.StatusCode
	if st != http.StatusOK {
		clog.WithField("status", st).Info("Process schedule response status is not OK.")
		return errors.New("Process schedule response status is not OK")
	}
	return nil
}

func (c *connector) EnableTimetableSync() error {
	u := new(url.URL)
	*u = *c.baseURL
	u.Path += syncPath
	enableURL := u.String()
	clog := log.WithFields(log.Fields{"logger": "init.annongen", "method": "EnableTimetableSync"})
	clog.WithField("url", enableURL).Info("Enabling annongen synchronization.")
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(enableURL); err != nil {
		clog.WithField("url", enableURL).WithError(err).WithField("url", u.RequestURI()).Info("Error while calling annongen synchronization URL.")
		return err
	}
	st := resp.StatusCode
	if st != http.StatusOK {
		clog.WithField("status", st).Info("Enable timetable sync response status is not OK.")
		return errors.New("Enable timetable sync response status is not OK")
	}
	return nil
}
