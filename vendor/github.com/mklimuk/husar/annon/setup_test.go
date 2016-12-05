package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/i18n"
	"github.com/mklimuk/husar/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SetupTestSuite struct {
	suite.Suite
	prod    Producer
	gen     Generator
	tpl     CatalogTemplate
	en      CatalogTemplate
	refTime time.Time
}

func (suite *SetupTestSuite) SetupSuite() {
	suite.prod = NewSetup()
	tpls := Parse([]string{"./templates.yml", "./templates_en.yml"})
	for _, t := range tpls {
		if t.ID == "setup" {
			t.Templates.Parse("test templates")
			suite.tpl = t
		}
		if t.ID == "setup_en" {
			t.Templates.Parse("en test templates")
			suite.en = t
		}
	}
	suite.gen = NewGenerator(Producers{suite.prod.Name(): &(suite.prod)}, prios)
	suite.refTime = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
}

func (suite *SetupTestSuite) SetupTest() {

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *SetupTestSuite) TestPassThrough() {
	t := test.GetTrain("passThrough")
	te, _, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), te)
}

func (suite *SetupTestSuite) TestStarting() {
	t := test.GetTrain("starting")
	te, update, _ := suite.prod.Required(t, nil, &suite.refTime)

	assert.True(suite.T(), te)
	assert.True(suite.T(), update)
	announcement, err := suite.gen.Generate(t, &suite.refTime, &suite.tpl, &suite.prod, i18n.GetLang(""))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), time.Date(2016, time.April, 10, 9, 5, 30, 0, config.Timezone), announcement.Time[0])
	assert.Equal(suite.T(), "Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 stoi na torze 1 przy peronie II. Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 stoi na torze 1 przy peronie II. Planowy odjazd pociągu o godzinie 09:08.", announcement.Text.HumanText)
	assert.Equal(suite.T(), `<s xml:lang="pl">Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 stoi na torze pierwszym przy peronie drugim. Pociąg Osobowy Kolei Mazowieckich RADOMIAK do stacji Station3 przez stacje: Station4, Station5 stoi na torze pierwszym przy peronie drugim. Planowy odjazd pociągu o godzinie 09:08.</s>`, announcement.Text.TtsText)
	assert.Equal(suite.T(), `Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, stoi na torze <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>. Pociąg <span class="tpl-param" data-name="category">Osobowy</span> <span class="tpl-param" data-name="carrier">Kolei Mazowieckich</span> <span class="tpl-param" data-name="name">RADOMIAK</span> do stacji <span class="tpl-param" data-name="category">Station3</span>, przez stacje: <span class="tpl-param" data-name="by">Station4, Station5</span>, stoi na torze <span class="tpl-param" data-name="track">1</span> przy peronie <span class="tpl-param" data-name="platform">II</span>.<br/>Planowy odjazd pociągu o godzinie <span class="tpl-param" data-name="departure">09:08</span>.`, announcement.Text.HTMLText)
	assert.Equal(suite.T(), false, announcement.Autoplay)
}

func (suite *SetupTestSuite) TestNoUpdate() {
	t := test.GetTrain("passThrough")
	_, update, _ := suite.prod.Required(t, t, &suite.refTime)
	assert.False(suite.T(), update)
	_, update, _ = suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), update)
}

func (suite *SetupTestSuite) TestEnding() {
	t := test.GetTrain("ending")
	te, _, _ := suite.prod.Required(t, nil, &suite.refTime)
	assert.False(suite.T(), te)

}

func TestSetupTestSuite(t *testing.T) {
	suite.Run(t, new(SetupTestSuite))
}
