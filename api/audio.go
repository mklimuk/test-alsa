package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mklimuk/husar/errors"
	"github.com/mklimuk/husar/rest"
	"github.com/mklimuk/test-alsa/audio"
	"github.com/mklimuk/websocket"

	log "github.com/Sirupsen/logrus"
)

type playAPI struct {
	a       audio.Playback
	factory websocket.ConnectionFactory
}

//NewPlaybackAPI is the playback API constructor
func NewPlaybackAPI(a audio.Playback, factory websocket.ConnectionFactory) rest.API {
	p := playAPI{a, factory}
	return rest.API(&p)
}

func (p *playAPI) AddRoutes(router *gin.Engine) {
	router.POST("/audio/play", p.play)
}

// play establishes a websocket connection and writes binary data to the playback device
func (p *playAPI) play(ctx *gin.Context) {
	log.WithFields(log.Fields{"logger": "audio-endpoint.api", "method": "play"}).
		Info("Establishing playback connection")
	req := new(audio.StreamContext)
	var err error
	if err = ctx.BindJSON(req); err != nil {
		log.WithFields(log.Fields{"logger": "audio-endpoint.api", "method": "play"}).Warn("Could not parse request")
		ctx.JSON(http.StatusBadRequest, errors.GetCtx(err))
		return
	}
	var c websocket.Connection
	if c, err = p.factory.UpgradeConnection(ctx.Writer, ctx.Request, nil); err != nil {

	}
	if err = p.a.PlayFromWsConnection(c, req); err != nil {

	}
}
