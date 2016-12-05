package annongen

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

const (
	healthPath string = "/health"
)

//CheckHealth calls annongen's health API to check whether the service is up and running
func (c *connector) CheckHealth() (bool, error) {
	clog := log.WithFields(log.Fields{"logger": "init.annongen", "method": "CheckHealth"})
	u := new(url.URL)
	*u = *c.baseURL
	u.Path += healthPath
	healthURL := u.String()
	client := &http.Client{}
	var resp *http.Response
	var err error
	if resp, err = client.Get(healthURL); err != nil {
		clog.WithError(err).WithField("url", healthURL).Error("Could not perform check health request.")
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}
