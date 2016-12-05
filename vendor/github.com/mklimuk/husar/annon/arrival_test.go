package annon

import (
	"testing"
	"time"

	"golang.org/x/text/language"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/i18n"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ArrivalTestSuite struct {
	suite.Suite
	prod    Producer
	tpl     CatalogTemplate
	en      CatalogTemplate
	gen     Generator
	refTime time.Time
}

func (suite *ArrivalTestSuite) SetupSuite() {
	suite.prod = NewArrival()
	tpls := Parse([]string{"./templates.yml", "./templates_en.yml"})
	for _, t := range tpls {
		if t.ID == "arrival" {
			t.Templates.Parse("test templates")
			suite.tpl = t
		}
		if t.ID == "arrival_en" {
			t.Templates.Parse("en test templates")
			suite.en = t
		}
	}
	suite.gen = NewGenerator(Producers{suite.prod.Name(): &(suite.prod)}, prios)
	suite.refTime = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
}

func (suite *ArrivalTestSuite) SetupTest() {

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ArrivalTestSuite) TestPassThrough() {
	t := test.GetTrain("passThrough")
	te, update, err := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.DefaultLang)
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), "Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Planowy odjazd pociągu o godzinie <span class="tpl-param" data-name="departure">15:23</span>.`, announcement.Text.HTMLText)
	t.Delay = 10
	announcement, err = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.DefaultLang)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Opóźniony pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Opóźniony pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Opóźniony pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Opóźniony pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Planowy odjazd pociągu o godzinie <span class="tpl-param" data-name="departure">15:23</span>.`, announcement.Text.HTMLText)
	_, err = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, language.English)
	assert.NoError(suite.T(), err)
	t.Services = &[]train.Service{
		train.Service{ID: "REZO", Name: "rezerwacja obowiązkowa"},
		train.Service{ID: "KOND", Name: "przewóz przesyłek konduktorskich", Carriage: []string{"12"}},
	}
	announcement, err = suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Opóźniony pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Opóźniony pociąg TLK LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor 1 przy peronie II. Pociąg jest objęty obowiązkową rezerwacją miejsc. Przesyłki konduktorskie przyjmuje i wydaje kierownik pociągu w wagonie numer 12. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Opóźniony pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Opóźniony pociąg <say-as interpret-as="cardinal">TLK</say-as> LUNA ze stacji Station2 do stacji Station3 przez stacje: Station6, Station7 wjedzie na tor pierwszy przy peronie drugim. Pociąg jest objęty obowiązkową rezerwacją miejsc. Przesyłki konduktorskie przyjmuje i wydaje kierownik pociągu w wagonie numer <say-as interpret-as="cardinal">12</say-as>. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu. Planowy odjazd pociągu o godzinie 15:23.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Opóźniony pociąg <span class="tpl-param" data-name="category">TLK</span> <span class="tpl-param" data-name="name">LUNA</span> ze stacji <span class="tpl-param" data-name="from">Station2</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station6, Station7</span>, wjedzie na tor <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Pociąg jest objęty obowiązkową rezerwacją miejsc.<br/>Przesyłki konduktorskie przyjmuje i wydaje kierownik pociągu w wagonie numer 12.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.<br/>Planowy odjazd pociągu o godzinie <span class="tpl-param" data-name="departure">15:23</span>.`, announcement.Text.HTMLText)
}

func (suite *ArrivalTestSuite) TestNoUpdate() {
	t := test.GetTrain("passThrough")
	_, update, _ := suite.prod.Required(t, t, &suite.refTime)
	assert.False(suite.T(), update)
	_, update, _ = suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), update)
}

func (suite *ArrivalTestSuite) TestStarting() {
	t := test.GetTrain("starting")
	te, update, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), te)
	assert.False(suite.T(), update)
}

func (suite *ArrivalTestSuite) TestEnding() {
	t := test.GetTrain("ending")
	te, update, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Pociąg Osobowy Kolei Mazowieckich ze stacji Station2 wjedzie na tor 8 przy peronie III. Pociąg Osobowy Kolei Mazowieckich ze stacji Station2 wjedzie na tor 8 przy peronie III. Pociąg kończy bieg. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Pociąg <say-as interpret-as="cardinal">Osobowy</say-as> Kolei Mazowieckich ze stacji Station2 wjedzie na tor ósmy przy peronie trzecim. Pociąg <say-as interpret-as="cardinal">Osobowy</say-as> Kolei Mazowieckich ze stacji Station2 wjedzie na tor ósmy przy peronie trzecim. Pociąg kończy bieg. Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> ze stacji <span class="tpl-param" data-name="from">Station2</span>, wjedzie na tor <span class="tpl-param" data-name="track">8</span> przy peronie <span class="tpl-param" data-name="platform">III</span>. Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> ze stacji <span class="tpl-param" data-name="from">Station2</span>, wjedzie na tor <span class="tpl-param" data-name="track">8</span> przy peronie <span class="tpl-param" data-name="platform">III</span>.<br/>Pociąg kończy bieg.<br/>Prosimy zachować ostrożność i nie zbliżać się do krawędzi peronu.`, announcement.Text.HTMLText)
}

func (suite *ArrivalTestSuite) TestManual() {
	t := test.GetTrain("passThrough")
	t.Settings = &train.Settings{Mode: train.Manual, Arrival: &train.TimetableEvent{
		Track:    15,
		Platform: "IV",
	}}
	p, _, _, _ := suite.prod.BuildParams(t, i18n.GetLang(""))
	assert.Equal(suite.T(), "IV", (*p)["Platform"].(string))
	assert.Equal(suite.T(), 15, (*p)["Track"].(int))
}

func (suite *ArrivalTestSuite) TestDelayedManual() {
	t := test.GetTrain("passThrough")
	t.Settings = &train.Settings{Mode: train.Manual, Delay: 10, Arrival: &train.TimetableEvent{
		Track:    15,
		Platform: "IV",
		Time:     time.Date(2016, time.April, 10, 15, 31, 0, 0, config.Timezone),
	}}
	ti, _, _ := suite.prod.GetTime(t, &suite.refTime)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 15, 26, 0, 0, config.Timezone), ti[0])
}

func (suite *ArrivalTestSuite) TestDelayedAuto() {
	t := test.GetTrain("passThrough")
	t.Delay = 10
	ti, f, l := suite.prod.GetTime(t, &suite.refTime)
	res := time.Date(2016, time.April, 10, 15, 26, 0, 0, config.Timezone)
	assert.Equal(suite.T(), res, ti[0])
	assert.Equal(suite.T(), res, f)
	assert.Equal(suite.T(), res, l)
}

func TestArrivalTestSuite(t *testing.T) {
	suite.Run(t, new(ArrivalTestSuite))
}
