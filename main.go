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
	p := audio.New(conf, alsa.NewPlaybackDevice(params), "/etc/husar/dong.wav")
	f := websocket.NewFactory()

	clog.Info("Initializing REST router...")
	z := api.NewPlaybackAPI(p, f)

	router := gin.New()
	z.AddRoutes(router)

	clog.Fatal(http.ListenAndServe(":8081", router))

}
