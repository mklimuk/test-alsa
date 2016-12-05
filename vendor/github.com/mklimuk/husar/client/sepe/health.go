package sepe

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

const (
	healthPath string = "/manage/health"
)

//CheckHealth calls SEPE's health API to check whether the service is up and running
func (c *connector) CheckHealth() (bool, error) {
	u := new(url.URL)
	*u = *c.monitorBaseURL
	u.Path += healthPath
	healthURL := u.String()
	clog := log.WithFields(log.Fields{"logger": "init.sepe", "method": "CheckHealth"})
	clog.WithField("url", healthURL).Info("Checking SEPE interface health.")
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(healthURL); err != nil {
		clog.WithError(err).WithField("url", healthURL).Info("Error while calling SEPE health URL.")
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		clog.WithField("status", resp.StatusCode).Info("Status message not OK.")
		return false, err
	}
	return true, nil
}
