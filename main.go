package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	alsa "github.com/mklimuk/test-alsa/alsa"
	"github.com/mklimuk/test-alsa/api"
	"github.com/mklimuk/test-alsa/audio"
	"github.com/mklimuk/test-alsa/config"
	"github.com/mklimuk/websocket"
)

const (
	defaultLogLevel string = "warn"
)

func main() {
	clog := log.WithFields(log.Fields{"logger": "mic-receiver.main"})

	log.SetLevel(log.InfoLevel)

	conf := &config.AudioConf{DeviceBuffer: 4096, ReadBuffer: 8192, PeriodFrames: 2048, Periods: 2}
	params := &audio.BufferParams{BufferFrames: conf.DeviceBuffer, PeriodFrames: conf.PeriodFrames, Periods: conf.Periods}


	var err error
	var dev audio.PlaybackDevice
	if dev, err = alsa.NewPlaybackDevice(params); err != nil {
		panic(err)
	}

	p := audio.New(conf, , "/etc/husar/dong.wav")
	f := websocket.NewFactory()

	clog.Info("Initializing REST router...")
	if z := api.NewPlaybackAPI(p, f)

	router := gin.New()
	z.AddRoutes(router)

	clog.Fatal(http.ListenAndServe(":8081", router))

}
