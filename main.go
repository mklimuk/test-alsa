package main

import (
	"flag"
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

	//read config path and log level from runtime flags
	//var configPath string
	//flag.StringVar(&configPath, "config", defaultConfigPath, "Path to yaml configuration file")
	level := flag.String("log", "info", "Log level")
	flag.Parse()

	var l log.Level
	var err error
	if l, err = log.ParseLevel(*level); err != nil {
		clog.WithField("level", level).Warn("Invalid log level. Fallback to info.")
		l = log.InfoLevel
	}
	log.SetLevel(l)

	conf := &config.AudioConf{DeviceBuffer: 2048, ReadBuffer: 2048, PeriodFrames: 1024, Periods: 2}
	params := &audio.BufferParams{BufferFrames: conf.DeviceBuffer, PeriodFrames: conf.PeriodFrames, Periods: conf.Periods}

	var dev audio.PlaybackDevice
	if dev, err = alsa.NewPlaybackDevice(params); err != nil {
		panic(err)
	}

	p := audio.New(conf, dev, "/etc/husar/dong.wav")
	f := websocket.NewFactory()

	clog.Info("Initializing REST router...")
	z := api.NewPlaybackAPI(p, f)

	router := gin.New()
	z.AddRoutes(router)

	clog.Fatal(http.ListenAndServe(":8081", router))

}
