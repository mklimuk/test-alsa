package audio

import (
	"github.com/mklimuk/husar/device"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayOne(t *testing.T) {
	expectedID := "abcd"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.FormValue("id")
		assert.Equal(t, expectedID, id)
		w.Write([]byte("ok"))

	}))
	defer ts.Close()

	co := NewController()
	c := co.(*con)
	res, code, err := c.play(ts.URL, expectedID)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.Equal(t, 200, code)
}

func TestPlay(t *testing.T) {
	expectedID := "abcd"
	var endpointsHosts []*httptest.Server
	expectedRes := "ok"

	for i := 0; i < 5; i++ {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		}))
		defer s.Close()
		endpointsHosts = append(endpointsHosts, s)
	}
	var devs []*device.Device
	for _, h := range endpointsHosts {
		u := strings.Split(h.URL, ":")
		devs = append(devs, &device.Device{IP: u[1][2:], Port: u[2]})
	}

	co := NewController()
	c := co.(*con)

	resps, statuses, errs := c.Play(devs, expectedID, "testing", time.Now())
	for i := 0; i < len(resps); i++ {
		assert.Equal(t, resps[i], expectedRes)
		assert.Equal(t, statuses[i], 200)
	}
	assert.Equal(t, len(errs), 0)
}

func TestPlayAll(t *testing.T) {

	expectedRes := "ok"
	expectedDesc := "desc"
	expectedID := "abcd"
	var endpointsHosts []*httptest.Server

	for i := 0; i < 5; i++ {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, expectedID, r.FormValue("id"))
			w.Write([]byte("ok"))
		}))
		defer s.Close()
		endpointsHosts = append(endpointsHosts, s)
	}

	var devs []*device.Device
	for _, h := range endpointsHosts {
		u := strings.Split(h.URL, ":")
		devs = append(devs, &device.Device{IP: u[1][2:], Port: u[2]})
	}

	co := NewController()
	c := co.(*con)

	resps, codes, errs := c.Play(devs, expectedID, expectedDesc, time.Now())

	for i := 0; i < len(resps); i++ {
		assert.Equal(t, resps[i], expectedRes)
		assert.Equal(t, codes[i], 200)
	}
	assert.Equal(t, len(errs), 0)
}
