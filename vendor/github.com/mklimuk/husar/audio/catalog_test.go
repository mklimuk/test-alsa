package audio

import (
	shared "github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/tts"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	log "github.com/Sirupsen/logrus"
)

type AudioCatalogTestSuite struct {
	suite.Suite
}

func (suite *AudioCatalogTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *AudioCatalogTestSuite) TestGetID() {
	c := NewCatalog(os.TempDir(), nil)
	assert.NotNil(suite.T(), c.GetID("test"))
}

func (suite *AudioCatalogTestSuite) TestGenerate() {
	a := shared.TTSMock{}
	resp := []byte{1, 2, 3, 4}
	t := uuid.New().String()
	a.On("GetAudio", mock.Anything, t).Once().Return(resp, 10, nil)
	au := tts.Connector(&a)
	c := NewCatalog(os.TempDir(), &au)
	ID, exists, duration, approxDur, err := c.Generate(t)
	assert.NotNil(suite.T(), ID)
	assert.Equal(suite.T(), 10, duration)
	assert.False(suite.T(), exists)
	assert.False(suite.T(), approxDur)
	assert.NoError(suite.T(), err)
	ID, exists, duration, approxDur, err = c.Generate(t)
	assert.NotNil(suite.T(), ID)
	assert.Equal(suite.T(), 10, duration)
	assert.True(suite.T(), exists)
	assert.False(suite.T(), approxDur)
	assert.NoError(suite.T(), err)
	a.AssertExpectations(suite.T())
}

func (suite *AudioCatalogTestSuite) TestGenerateGet() {
	a := shared.TTSMock{}
	resp := []byte{1, 2, 3, 4}
	t := uuid.New().String()
	a.On("GetAudio", mock.Anything, t).Once().Return(resp, 10000, nil)
	au := tts.Connector(&a)
	c := NewCatalog(os.TempDir(), &au)
	ID, _, _, _, _ := c.Generate(t)
	data, _ := c.Get(ID)
	assert.Equal(suite.T(), resp, data)
}

func TestAudioCatalogTestSuite(t *testing.T) {
	suite.Run(t, new(AudioCatalogTestSuite))
}
