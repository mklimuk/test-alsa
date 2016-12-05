package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/i18n"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DepartureTestSuite struct {
	suite.Suite
	prod    Producer
	gen     Generator
	tpl     CatalogTemplate
	en      CatalogTemplate
	refTime time.Time
}

func (suite *DepartureTestSuite) SetupSuite() {
	suite.prod = NewDeparture()
	tpls := Parse([]string{"./templates.yml", "./templates_en.yml"})
	for _, t := range tpls {
		if t.ID == "departure" {
			t.Templates.Parse("test templates")
			suite.tpl = t
		}
		if t.ID == "departure_en" {
			t.Templates.Parse("en test templates")
			suite.en = t
		}
	}
	suite.gen = NewGenerator(Producers{suite.prod.Name(): &(suite.prod)}, prios)
	suite.refTime = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
}

func (suite *DepartureTestSuite) SetupTest() {

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *DepartureTestSuite) TestPassThrough() {
	t := test.GetTrain("passThrough")
	te, update, err := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Życzymy Państwu przyjemnej podróży.`, announcement.Text.HTMLText)
	t.Delay = 10
	announcement, err = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []time.Time{time.Date(2016, time.April, 10, 15, 32, 30, 0, config.Timezone)}, announcement.Time)
	assert.Equal(suite.T(), "Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Życzymy Państwu przyjemnej podróży.`, announcement.Text.HTMLText)
	t.Services = &[]train.Service{
		train.Service{ID: "REZO", Name: "rezerwacja obowiązkowa"},
		train.Service{ID: "KOND", Name: "przewóz przesyłek konduktorskich", Carriage: []string{"12"}},
	}
	announcement, err = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Opóźniony pociąg TLK LUNA do stacji Station3 przez stacje: Station6, Station7 odjedzie z toru pierwszego przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Życzymy Państwu przyjemnej podróży.`, announcement.Text.HTMLText)

}

func (suite *DepartureTestSuite) TestNoUpdate() {
	t := test.GetTrain("passThrough")
	_, update, _ := suite.prod.Required(t, t, &suite.refTime)
	assert.False(suite.T(), update)
	_, update, _ = suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), update)
}

func (suite *DepartureTestSuite) TestStarting() {
	t := test.GetTrain("starting")
	te, update, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 odjedzie z toru 1 przy peronie II. Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 odjedzie z toru 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 odjedzie z toru pierwszego przy peronie drugim. Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 odjedzie z toru pierwszego przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Życzymy Państwu przyjemnej podróży.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, odjedzie z toru <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Życzymy Państwu przyjemnej podróży.`, announcement.Text.HTMLText)

}

func (suite *DepartureTestSuite) TestEnding() {
	t := test.GetTrain("ending")
	te, _, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), te)
}

func (suite *DepartureTestSuite) TestDelayedManual() {
	t := test.GetTrain("passThrough")
	t.Settings = &train.Settings{Mode: train.Manual, Delay: 10, Departure: &train.TimetableEvent{
		Track:    15,
		Platform: "IV",
		Time:     time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone),
	}}
	ti, _, _ := suite.prod.GetTime(t, &suite.refTime)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 32, 30, 0, config.Timezone), ti[0])
}

func (suite *DepartureTestSuite) TestDelayedAuto() {
	t := test.GetTrain("passThrough")
	t.Delay = 10
	ti, f, l := suite.prod.GetTime(t, &suite.refTime)
	res := time.Date(2016, time.April, 10, 15, 32, 30, 0, config.Timezone)
	assert.Equal(suite.T(), res, ti[0])
	assert.Equal(suite.T(), res, f)
	assert.Equal(suite.T(), res, l)
}

func TestDepartureTestSuite(t *testing.T) {
	suite.Run(t, new(DepartureTestSuite))
}
