package api

import (
	"github.com/gin-gonic/gin"
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
	router.GET("/audio/play", p.play)
}

// play establishes a websocket connection and writes binary data to the playback device
func (p *playAPI) play(ctx *gin.Context) {
	defer rest.ErrorHandler(ctx)
	log.WithFields(log.Fields{"logger": "audio-endpoint.api", "method": "play"}).
		Info("Establishing playback connection")
	var c websocket.Connection
	var err error
	if c, err = p.factory.UpgradeConnection(ctx.Writer, ctx.Request, nil); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	if err = p.a.PlayFromWsConnection(c); err != nil {
		c.Close("")
	}
}
