package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/mklimuk/test-alsa/audio"
	"github.com/mklimuk/websocket"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PlaybackAPITestSuite struct {
	suite.Suite
	router *gin.Engine
	serv   *httptest.Server
	a      playAPI
}

func (suite *PlaybackAPITestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	w := websocket.FactoryMock{}
	p := audio.PlaybackMock{}
	suite.a = playAPI{&p, &w}
	suite.router = gin.New()
	suite.a.AddRoutes(suite.router)
	suite.serv = httptest.NewServer(suite.router)
}

func (suite *PlaybackAPITestSuite) TearDownSuite() {
	suite.serv.Close()
}

func (suite *PlaybackAPITestSuite) TestRequestError() {
	f := websocket.FactoryMock{}
	f.On("UpgradeConnection", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mock error")).Once()
	suite.a.factory = &f
	res, err := http.Get(fmt.Sprintf("%s/%s", suite.serv.URL, "audio/play"))
	a := assert.New(suite.T())
	a.NoError(err)
	a.Equal(400, res.StatusCode)
}

func (suite *PlaybackAPITestSuite) TestVersion() {

}

func TestPlaybackAPITestSuite(t *testing.T) {
	suite.Run(t, new(PlaybackAPITestSuite))
}
