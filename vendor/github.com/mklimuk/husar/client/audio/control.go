package audio

import (
	"fmt"
	"github.com/mklimuk/husar/device"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func (c *con) SetVolume(devices []*device.Device, volume int) ([]string, []int, []error) {
	var errs []error
	var responses []string
	var statusesCode []int
	var wg sync.WaitGroup
	for _, d := range devices {
		wg.Add(1)
		go func(IP string, port string, vol int) {
			response, statusCode, err := c.setVol(fmt.Sprintf("%s:%s/amp/control/volume/%d", IP, port, vol))
			responses = append(responses, response)
			statusesCode = append(statusesCode, statusCode)
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(d.IP, d.Port, volume)
	}
	wg.Wait()
	return responses, statusesCode, errs

}

func (c *con) setVol(apiURL string) (string, int, error) {
	if apiURL[0:7] != "http://" && apiURL[0:8] != "https://" {
		apiURL = "http://" + apiURL
	}
	r, err := http.NewRequest("PUT", apiURL, nil) // <-- URL-encoded payload
	if err != nil {
		log.WithFields(log.Fields{"logger": "audio.control", "method": "setVol", "device": apiURL}).
			WithError(err).Error("Could not build HTTP request.")
		return "", 0, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		log.WithFields(log.Fields{"logger": "audio.control", "method": "setVol", "device": apiURL}).
			WithError(err).Error("HTTP request has failed.")
		if resp != nil {
			return resp.Status, resp.StatusCode, err
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"logger": "audio.control", "method": "setVol", "device": apiURL}).
			WithError(err).Error("Could not parse HTTP response.")
		return resp.Status, resp.StatusCode, err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "audio.control", "method": "setVol",
			"device": apiURL, "status": resp.StatusCode, "response": string(respBody)}).
			Debug("Response from device device.")
	}
	return string(respBody), resp.StatusCode, err
}
