package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/i18n"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DelayTestSuite struct {
	suite.Suite
	prod    Producer
	gen     Generator
	tpl     CatalogTemplate
	refTime time.Time
}

func (suite *DelayTestSuite) SetupSuite() {
	suite.prod = NewDelay()
	path := "./templates.yml"
	tpls := Parse([]string{path})
	for _, t := range tpls {
		if t.Type == suite.prod.Name() {
			t.Templates.Parse("test templates")
			suite.tpl = t
		}
	}
	suite.gen = NewGenerator(Producers{suite.prod.Name(): &(suite.prod)}, prios)
	suite.refTime = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
}

func (suite *DelayTestSuite) SetupTest() {

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *DelayTestSuite) TestNoDelay() {
	t := test.GetTrain("passThrough")
	t.Delay = 0
	te, _, err := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), te)
	assert.Nil(suite.T(), err)
}

func (suite *DelayTestSuite) TestRequired() {
	t := test.GetTrain("passThrough")
	t.Delay = 10
	te, update, err := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	assert.Nil(suite.T(), err)
	te, update, _ = suite.prod.Required(t, t, &suite.refTime)
	assert.True(suite.T(), te)
	assert.False(suite.T(), update)
	changedDelay := train.Train{
		Delay: 20,
	}
	te, update, _ = suite.prod.Required(t, &changedDelay, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
}

func (suite *DelayTestSuite) TestPassThrough() {
	t := test.GetTrain("passThrough")
	t.Delay = 10
	delayed := t.Arrival.Time.Add(time.Duration(t.Delay) * time.Minute)
	t.FirstLiveEvent = &delayed
	humanText10 := "Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7, planowy przyjazd o godzinie 15:21, przyjedzie z opóźnieniem około 10 minut. Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7, planowy przyjazd o godzinie 15:21, przyjedzie z opóźnieniem około 10 minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy."
	humanText20 := "Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7, planowy przyjazd o godzinie 15:21, przyjedzie z opóźnieniem około 20 minut. Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7, planowy przyjazd o godzinie 15:21, przyjedzie z opóźnieniem około 20 minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy."
	ttsText10 := `<s xml:lang="pl">Pociąg <emphasis level="strong">TLK</emphasis> <emphasis level="strong">LUNA</emphasis>, ze stacji: <emphasis level="strong">Station2</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station6, Station7</emphasis>, planowy przyjazd o godzinie 15:21 , przyjedzie z opóźnieniem około <emphasis level="strong">10</emphasis> minut. Pociąg <emphasis level="strong">TLK</emphasis> <emphasis level="strong">LUNA</emphasis>, ze stacji: <emphasis level="strong">Station2</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station6, Station7</emphasis>, planowy przyjazd o godzinie 15:21 , przyjedzie z opóźnieniem około <emphasis level="strong">10</emphasis> minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy.</s>`
	ttsText20 := `<s xml:lang="pl">Pociąg <emphasis level="strong">TLK</emphasis> <emphasis level="strong">LUNA</emphasis>, ze stacji: <emphasis level="strong">Station2</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station6, Station7</emphasis>, planowy przyjazd o godzinie 15:21 , przyjedzie z opóźnieniem około <emphasis level="strong">20</emphasis> minut. Pociąg <emphasis level="strong">TLK</emphasis> <emphasis level="strong">LUNA</emphasis>, ze stacji: <emphasis level="strong">Station2</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station6, Station7</emphasis>, planowy przyjazd o godzinie 15:21 , przyjedzie z opóźnieniem około <emphasis level="strong">20</emphasis> minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy.</s>`
	htmlText10 := `Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, planowy przyjazd o godzinie <span class="tpl-param" data-name="arrival">15:21</span> przyjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">10</span> minut. Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, planowy przyjazd o godzinie <span class="tpl-param" data-name="arrival">15:21</span> przyjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">10</span> minut.<br/>Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty.<br/>Za opóźnienie pociągu przepraszamy.`
	htmlText20 := `Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, planowy przyjazd o godzinie <span class="tpl-param" data-name="arrival">15:21</span> przyjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">20</span> minut. Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, planowy przyjazd o godzinie <span class="tpl-param" data-name="arrival">15:21</span> przyjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">20</span> minut.<br/>Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty.<br/>Za opóźnienie pociągu przepraszamy.`
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	require.Len(suite.T(), announcement.Time, 3, "There should be 3 announcement.")
	assert.Equal(suite.T(), []time.Time{time.Date(2016, time.April, 10, 15, 11, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 15, 16, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)}, announcement.Time)
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), humanText10, announcement.Text.HumanText)
	assert.Equal(suite.T(), ttsText10, announcement.Text.TtsText)
	assert.Equal(suite.T(), htmlText10, announcement.Text.HTMLText)
	t.Delay = 20
	delayed = t.Arrival.Time.Add(time.Duration(t.Delay) * time.Minute)
	t.FirstLiveEvent = &delayed
	announcement, _ = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Equal(suite.T(), humanText20, announcement.Text.HumanText)
	assert.Equal(suite.T(), ttsText20, announcement.Text.TtsText)
	assert.Equal(suite.T(), htmlText20, announcement.Text.HTMLText)
}

func (suite *DelayTestSuite) TestCleanup() {
	t := test.GetTrain("passThrough")
	old := new(train.Train)
	train.Copy(t, old)
	old.Delay = 10
	is, update, _ := suite.prod.Required(t, old, &suite.refTime)
	assert.False(suite.T(), is)
	assert.True(suite.T(), update)
	old.Delay = 0
	is, update, _ = suite.prod.Required(t, old, &suite.refTime)
	assert.False(suite.T(), is)
	assert.False(suite.T(), update)
}

func (suite *DelayTestSuite) TestAfterArrival() {
	t := test.GetTrain("passThrough")
	t.Delay = 120
	t.Arrival.Time = time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	t.Departure.Time = time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone)
	delayed := t.Arrival.Time.Add(time.Duration(t.Delay) * time.Minute)
	t.FirstLiveEvent = &delayed
	ref := time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone)
	announcement, err := suite.gen.Generate(t, &ref, &suite.tpl, &suite.prod, i18n.GetLang(""))
	require.Len(suite.T(), announcement.Time, 4, "There should be 4 announcement.")
	assert.Equal(suite.T(), []time.Time{time.Date(2016, time.April, 10, 15, 25, 30, 0, config.Timezone),
		time.Date(2016, time.April, 10, 15, 55, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 16, 25, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 16, 55, 0, 0, config.Timezone),
	}, announcement.Time)
	assert.Equal(suite.T(), err, nil)
}

func (suite *DelayTestSuite) TestStarting() {
	t := test.GetTrain("starting")
	t.Delay = 15
	delayed := t.Departure.Time.Add(time.Duration(t.Delay) * time.Minute)
	t.FirstLiveEvent = &delayed
	humanText := "Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5, planowy odjazd o godzinie 09:08, odjedzie z opóźnieniem około 15 minut. Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5, planowy odjazd o godzinie 09:08, odjedzie z opóźnieniem około 15 minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy."
	ttsText := `<s xml:lang="pl">Pociąg <emphasis level="strong">Osobowy</emphasis> <emphasis level="strong">Kolei Mazowieckich</emphasis> <emphasis level="strong">RADOMIAK</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station4, Station5</emphasis>, planowy odjazd o godzinie 09:08 , odjedzie z opóźnieniem około <emphasis level="strong">15</emphasis> minut. Pociąg <emphasis level="strong">Osobowy</emphasis> <emphasis level="strong">Kolei Mazowieckich</emphasis> <emphasis level="strong">RADOMIAK</emphasis>, do stacji: <emphasis level="strong">Station3</emphasis>, przez stacje: <emphasis level="strong">Station4, Station5</emphasis>, planowy odjazd o godzinie 09:08 , odjedzie z opóźnieniem około <emphasis level="strong">15</emphasis> minut. Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty. Za opóźnienie pociągu przepraszamy.</s>`
	htmlText := `Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, planowy odjazd o godzinie <span class="tpl-param" data-name="arrival">09:08</span> odjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">15</span> minut. Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, planowy odjazd o godzinie <span class="tpl-param" data-name="arrival">09:08</span> odjedzie z opóźnieniem około <span class="tpl-param" data-name="delay">15</span> minut.<br/>Opóźnienie może ulec zmianie. Prosimy o zwracanie uwagi na komunikaty.<br/>Za opóźnienie pociągu przepraszamy.`
	te, update, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), humanText, announcement.Text.HumanText)
	assert.Equal(suite.T(), ttsText, announcement.Text.TtsText)
	assert.Equal(suite.T(), htmlText, announcement.Text.HTMLText)

}

func (suite *DelayTestSuite) TestManual() {
	t := test.GetTrain("passThrough")
	t.Delay = 10
	t.Settings = &train.Settings{Delay: 20, Mode: train.Manual}
	t.Settings.Arrival = &train.TimetableEvent{Time: t.Arrival.Time.Add(time.Duration(t.Delay) * time.Minute)}
	t.Settings.Departure = &train.TimetableEvent{Time: t.Departure.Time.Add(time.Duration(t.Delay) * time.Minute)}
	ref := time.Date(2016, time.April, 10, 15, 10, 0, 0, config.Timezone)
	ti, first, last := suite.prod.GetTime(t, &ref)
	require.Len(suite.T(), ti, 3, "There should be 3 announcement.")
	assert.Equal(suite.T(), []time.Time{time.Date(2016, time.April, 10, 15, 11, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 15, 16, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
	}, ti)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 11, 0, 0, config.Timezone), first)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone), last)
	p, _, _, _ := suite.prod.BuildParams(t, i18n.GetLang(""))
	assert.Equal(suite.T(), 20, (*p)["Delay"].(int))
	t.Settings.Mode = train.Auto
	p, _, _, _ = suite.prod.BuildParams(t, i18n.GetLang(""))
	assert.Equal(suite.T(), 10, (*p)["Delay"].(int))
}
func (suite *DelayTestSuite) TestEnding() {
	t := test.GetTrain("ending")
	test, _, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), test)
}

func TestDelayTestSuite(t *testing.T) {
	suite.Run(t, new(DelayTestSuite))
}
