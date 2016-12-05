package annon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	html "html/template"
	txt "text/template"
)

const tplTxt = "tor {{.TrackTxt}} przy peronie {{.Platform}}"
const tplHtm = "tor <span>{{.TrackTxt}}</span> przy peronie <span>{{.Platform}}</span>"
const expectTxt = "tor 1 przy peronie II"
const expectHtm = "tor <span>1</span> przy peronie <span>II</span>"
const expectTts = "tor pierwszy przy peronie drugim"

type TemplatesTestSuite struct {
	suite.Suite
	txt *txt.Template
	htm *html.Template
}

func (suite *TemplatesTestSuite) SetupSuite() {
	suite.txt = txt.Must(txt.New("test_txt").Parse(tplTxt))
	suite.htm = html.Must(html.New("test_htm").Parse(tplHtm))
}

func (suite *TemplatesTestSuite) TestTemplateParser() {
	tps := &Templates{human: suite.txt, tts: suite.txt, html: suite.htm}
	params := TemplateParams{"Platform": "II", "Track": 1, "TrackTxt": "1"}
	ttsParams := TemplateParams{"Platform": "drugim", "Track": 1, "TrackTxt": "pierwszy"}
	human, tts, html, err := ExecuteTemplates(tps, params, ttsParams)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), expectTts, tts)
	assert.Equal(suite.T(), expectTxt, human)
	assert.Equal(suite.T(), expectHtm, html)
}

func (suite *TemplatesTestSuite) TestMapParamsParser() {
	tps := &Templates{human: suite.txt, tts: suite.txt, html: suite.htm}
	params := TemplateParams{
		"Platform": "II", "Track": "1", "TrackTxt": "1"}
	human, tts, html, err := ExecuteTemplates(tps, params, params)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), expectTxt, tts)
	assert.Equal(suite.T(), expectTxt, human)
	assert.Equal(suite.T(), expectHtm, html)
}

func (suite *TemplatesTestSuite) TestParseMethod() {
	tps := &Templates{Human: "human", Tts: "tts", HTML: "html"}
	tps.Parse("test")
	assert.NotNil(suite.T(), tps.human)
	assert.NotNil(suite.T(), tps.tts)
	assert.NotNil(suite.T(), tps.html)
}

func TestTemplates(t *testing.T) {
	suite.Run(t, new(TemplatesTestSuite))
}
