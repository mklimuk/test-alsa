package audio

import (
	"bytes"
	"fmt"
	"github.com/mklimuk/husar/device"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func (c *con) Upload(audio []byte, title string, devices []*device.Device, id string, concurrent int) ([]string, []int, []error) {
	var errs []error
	var responses []string
	var statusesCode []int
	var wg sync.WaitGroup
	limit := make(chan int, concurrent)
	for _, d := range devices {
		wg.Add(1)
		limit <- 1
		go func(url string, port string) {
			uri := fmt.Sprintf("%s:%s/audio", url, port)
			response, statusCode, err := c.upload(audio, title, uri, id)
			responses = append(responses, response)
			statusesCode = append(statusesCode, statusCode)
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
			<-limit
		}(d.IP, d.Port)
	}
	wg.Wait()
	close(limit)
	return responses, statusesCode, errs
}

// upload sends audio file to the given endpoint
func (c *con) upload(audio []byte, title string, IP string, id string) (string, int, error) {
	// TODO check why we don't control the uri
	if IP[0:7] != "http://" && IP[0:8] != "https://" {
		IP = "http://" + IP
	}
	req, boundary, err := c.buildUploadRequest(IP, title, id, audio)
	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.upload", "method": "upload", "id": id, "device": IP}).
			WithError(err).Error("Could not build upload request.")
		return "", 0, err
	}
	req.Header.Add("Content-Type", boundary)
	req.Header.Add("Content-Type", "multipart/form-data")

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "device.audio.upload", "method": "upload", "id": id, "device": IP}).
			Debug("Sending data to device.")
	}
	resp, err := c.client.Do(req)

	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.upload", "method": "upload", "id": id, "device": IP}).
			WithError(err).Error("Could not send payload to device.")
		if resp != nil {
			return resp.Status, resp.StatusCode, err
		}
		return "", -1, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"logger": "device.audio.upload", "method": "upload", "id": id, "device": IP}).
			WithError(err).Error("Could not parse response body.")
		return resp.Status, resp.StatusCode, err
	}

	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "device.audio.upload", "method": "upload", "id": id,
			"device": IP, "status": resp.StatusCode, "response": string(respBody)}).
			Debug("Response from device device.")
	}
	return string(respBody), resp.StatusCode, nil
}

func (c *con) buildUploadRequest(uri string, title string, id string, audio []byte) (*http.Request, string, error) {

	var err error
	// create body to write form to
	body := &bytes.Buffer{}

	// create writer that will write into body
	writer := multipart.NewWriter(body)

	// add description field
	if err = writer.WriteField("description", title); err != nil {
		return nil, "", err
	}

	// add description field
	if err = writer.WriteField("id", id); err != nil {
		return nil, "", err
	}

	// add file field
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%s.ogg", id))
	if err != nil {
		return nil, "", err
	}

	// add file to file field
	part.Write(audio)
	boundary := writer.FormDataContentType()

	// close writer
	if err = writer.Close(); err != nil {
		return nil, "", err
	}

	// return request
	req, err := http.NewRequest("POST", uri, body)
	return req, boundary, err
}
