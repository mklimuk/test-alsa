package audio

import (
	"bytes"
	"fmt"
	"github.com/mklimuk/husar/device"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

func (c *con) Play(devices []*device.Device, id string, title string, timestamp time.Time) ([]string, []int, []error) {
	var errs []error
	var responses []string
	var statusesCode []int
	var wg sync.WaitGroup
	for _, d := range devices {
		wg.Add(1)
		go func(IP string, port string) {
			response, statusCode, err := c.play(fmt.Sprintf("%s:%s/audio", IP, port), id)
			responses = append(responses, response)
			statusesCode = append(statusesCode, statusCode)
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(d.IP, d.Port)
	}
	wg.Wait()
	return responses, statusesCode, errs

}

func (c *con) play(apiURL string, id string) (string, int, error) {
	data := url.Values{}
	data.Set("id", id)

	if apiURL[0:7] != "http://" && apiURL[0:8] != "https://" {
		apiURL = "http://" + apiURL
	}
	r, err := http.NewRequest("PUT", apiURL, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.play", "method": "play", "id": id, "device": apiURL}).
			WithError(err).Error("Could not build HTTP request.")
		return "", 0, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := c.client.Do(r)
	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.play", "method": "play", "id": id, "device": apiURL}).
			WithError(err).Error("HTTP request has failed.")
		if resp != nil {
			return resp.Status, resp.StatusCode, err
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.play", "method": "play", "id": id, "device": apiURL}).
			WithError(err).Error("Could not parse HTTP response.")
		if resp != nil {
			return resp.Status, resp.StatusCode, err
		}
		return "", -1, err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "device.audio.play", "method": "play", "id": id,
			"device": apiURL, "status": resp.StatusCode, "response": string(respBody)}).
			Debug("Response from device device.")
	}
	return string(respBody), resp.StatusCode, err
}
