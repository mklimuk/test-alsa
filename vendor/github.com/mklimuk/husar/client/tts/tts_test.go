package tts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAudio(t *testing.T) {

	requestedString := "test of tts"
	expectedURL := `/test%20of%20tts`
	expectedRes := []byte{1, 2, 3, 4}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.String(), expectedURL)
		w.Header().Add("Audio-Len", "10")
		w.Write(expectedRes)
	}))
	defer ts.Close()

	c := NewConnector(ts.URL + "/test%20of%20tts")
	res, len, err := c.GetAudio("test", requestedString)
	assert.NoError(t, err)
	assert.Equal(t, expectedRes, res)
	assert.Equal(t, 10, len)
}
