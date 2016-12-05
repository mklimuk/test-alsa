package tts

import (
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SpeechSuite struct {
	suite.Suite
	speech Speech
}

func (suite *SpeechSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	log.Debug(os.TempDir())
	suite.speech = NewSpeech(os.TempDir(), "./fake_synth.sh", "./fake_converter.sh", "-1", "test")
}

func (suite *SpeechSuite) TestLength() {
	len, err := audioLength("fake.wav")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 8, len)
}

func (suite *SpeechSuite) TestGenerate() {
	gr := GenerateRequest{
		ID:   "test",
		Text: "Hello world",
		Time: time.Now(),
	}
	path, len, err := suite.speech.GenerateSpeech(gr)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 8, len)
	assert.Contains(suite.T(), path, "test.ogg")
}

func TestSpeechSuite(t *testing.T) {
	suite.Run(t, new(SpeechSuite))
}
