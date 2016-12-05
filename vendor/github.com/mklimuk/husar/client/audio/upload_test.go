package audio

import (
	"github.com/mklimuk/husar/device"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadOne(t *testing.T) {
	expectedData := []byte{1, 2, 3, 4}
	expectedRes := "ok"
	expectedDesc := "desc"
	expectedID := "abcd"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(1024)
		assert.NoError(t, err)
		desc := r.FormValue("description")
		id := r.FormValue("id")
		file, header, err := r.FormFile("file")
		assert.NoError(t, err)
		defer file.Close()
		fiecontent, err := ioutil.ReadAll(file)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, fiecontent)
		assert.Equal(t, "abcd.ogg", header.Filename)
		assert.Equal(t, expectedDesc, desc)
		assert.Equal(t, expectedID, id)
		w.Write([]byte("ok"))

	}))
	defer ts.Close()

	co := NewController()
	c := co.(*con)

	res, code, err := c.upload(expectedData, expectedDesc, ts.URL, expectedID)
	assert.NoError(t, err)
	assert.Equal(t, expectedRes, res)
	assert.Equal(t, 200, code)
}

func TestUploadAll(t *testing.T) {
	expectedData := []byte{1, 2, 3, 4}
	expectedRes := "ok"
	expectedDesc := "desc"
	expectedID := "abcd"
	var endpointsHosts []*httptest.Server

	for i := 0; i < 5; i++ {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseMultipartForm(1024)
			assert.NoError(t, err)
			desc := r.FormValue("description")
			id := r.FormValue("id")
			file, header, err := r.FormFile("file")
			assert.NoError(t, err)
			defer file.Close()
			fiecontent, err := ioutil.ReadAll(file)
			assert.NoError(t, err)
			assert.Equal(t, expectedData, fiecontent)
			assert.Equal(t, "abcd.ogg", header.Filename)
			assert.Equal(t, expectedDesc, desc)
			assert.Equal(t, expectedID, id)
			w.Write([]byte("ok"))
		}))

		endpointsHosts = append(endpointsHosts, s)
		defer s.Close()
	}

	var devs []*device.Device
	for _, h := range endpointsHosts {
		u := strings.Split(h.URL, ":")
		devs = append(devs, &device.Device{IP: u[1][2:], Port: u[2]})
	}

	co := NewController()
	c := co.(*con)

	resps, codes, errs := c.Upload(expectedData, expectedDesc, devs, expectedID, 2)
	for i := 0; i < len(resps); i++ {
		assert.Equal(t, resps[i], expectedRes)
		assert.Equal(t, codes[i], 200)
	}
	assert.Equal(t, len(errs), 0)
}
