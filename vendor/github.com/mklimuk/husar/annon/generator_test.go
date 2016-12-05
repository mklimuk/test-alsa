package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
}

func (suite *GeneratorTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (suite *GeneratorTestSuite) SetupTest() {

}

// All methods that begin with "Test" are run as tests within a
// suite.

func (suite *GeneratorTestSuite) TestNotRequired() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	producer1 := &ProducerMock{}
	cat := Catalog(&CatalogMock{})
	newTrain := &train.Train{}
	oldTrain := &train.Train{}
	currentAnnons := &[]*Announcement{}
	producer1.On("Required", newTrain, oldTrain, mock.AnythingOfType("*time.Time")).Return(false, true, nil)
	producer1.On("Name").Return(Arrival)
	p1 := Producer(producer1)
	gen.addProducer(&p1)
	now := time.Now()
	annons, _, err := gen.ForTrain(newTrain, oldTrain, currentAnnons, &now, &cat)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), annons)
	producer1.AssertNotCalled(suite.T(), "ForTrain", newTrain, mock.AnythingOfType("*time.Time"))
}

func (suite *GeneratorTestSuite) TestRequiredNoUpdate() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	producer1 := &ProducerMock{}
	cat := Catalog(&CatalogMock{})
	newTrain := &train.Train{}
	oldTrain := &train.Train{}
	good := Announcement{
		ID:   "arrival",
		Type: Arrival,
	}
	bad := Announcement{
		ID:   "departure",
		Type: Departure,
	}
	currentAnnons := []*Announcement{
		&good,
		&bad,
	}
	producer1.On("Required", newTrain, oldTrain, mock.AnythingOfType("*time.Time")).Return(true, false, nil)
	producer1.On("Name").Return(Arrival)
	p1 := Producer(producer1)
	gen.addProducer(&p1)
	now := time.Now()
	annons, _, err := gen.ForTrain(newTrain, oldTrain, &currentAnnons, &now, &cat)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), annons, 1)
	producer1.AssertNotCalled(suite.T(), "ForTrain", newTrain, &now)
}

func (suite *GeneratorTestSuite) TestRequiredWithUpdate() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	gen.cache = make(map[string]*CatalogTemplate)
	now := time.Now()
	producer1 := &ProducerMock{}
	cat := &CatalogMock{}
	tem := &CatalogTemplate{
		Type:  Arrival,
		Title: "Test template",
	}
	cat.On("Get", mock.AnythingOfType("string")).Return(tem, nil)
	c := Catalog(cat)
	newTrain := &train.Train{}
	oldTrain := &train.Train{}
	params := &TemplateParams{}
	producer1.On("Required", newTrain, oldTrain, mock.AnythingOfType("*time.Time")).Return(true, true, nil)
	producer1.On("Name").Return(Arrival)
	producer1.On("GetTime", mock.AnythingOfType("*train.Train"), mock.AnythingOfType("*time.Time")).Return([]time.Time{
		time.Date(2016, time.April, 10, 23, 50, 0, 0, config.Timezone),
		time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone),
	}, time.Date(2016, time.April, 10, 23, 50, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone))
	producer1.On("BuildParams", mock.AnythingOfType("*train.Train")).Return(params, params, true, nil)
	p1 := Producer(producer1)
	gen.addProducer(&p1)
	annons, _, err := gen.ForTrain(newTrain, oldTrain, nil, &now, &c)
	producer1.AssertCalled(suite.T(), "BuildParams", mock.AnythingOfType("*train.Train"))
	assert.Len(suite.T(), annons, 1)
	assert.Len(suite.T(), annons[0].Time, 2)
	assert.Nil(suite.T(), err)
}

func (suite *GeneratorTestSuite) TestAggregation() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	gen.cache = make(map[string]*CatalogTemplate)
	now := time.Now()
	producer1 := &ProducerMock{}
	producer2 := &ProducerMock{}
	cat := &CatalogMock{}
	tem := &CatalogTemplate{
		ID:    "tarr",
		Type:  Arrival,
		Title: "Test template",
		Templates: Templates{
			Human: "human",
			Tts:   "tts",
			HTML:  "html",
		},
	}
	cat.On("Get", mock.AnythingOfType("string")).Return(tem, nil)
	c := Catalog(cat)
	newTrain := &train.Train{ID: "new"}
	oldTrain := &train.Train{ID: "old"}
	params := &TemplateParams{}
	producer1.On("Required", newTrain, oldTrain, mock.AnythingOfType("*time.Time")).Return(true, true, nil)
	producer1.On("GetTime", mock.AnythingOfType("*train.Train"), mock.AnythingOfType("*time.Time")).Return([]time.Time{
		time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone),
	}, time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone))
	producer1.On("BuildParams", mock.AnythingOfType("*train.Train")).Return(params, params, true, nil)
	producer1.On("Name").Return(Arrival)
	producer2.On("Required", newTrain, oldTrain, mock.AnythingOfType("*time.Time")).Return(true, true, nil)
	producer2.On("GetTime", mock.AnythingOfType("*train.Train"), mock.AnythingOfType("*time.Time")).Return([]time.Time{
		time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone),
	}, time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone), time.Date(2016, time.April, 10, 23, 59, 0, 0, config.Timezone))
	producer2.On("BuildParams", mock.AnythingOfType("*train.Train")).Return(params, params, true, nil)
	producer2.On("Name").Return(Departure)
	p1 := Producer(producer1)
	gen.addProducer(&p1)
	p2 := Producer(producer2)
	gen.addProducer(&p2)
	annons, _, err := gen.ForTrain(newTrain, oldTrain, nil, &now, &c)
	producer1.AssertCalled(suite.T(), "BuildParams", newTrain)
	producer2.AssertCalled(suite.T(), "BuildParams", newTrain)
	assert.Len(suite.T(), annons, 2)
	assert.Nil(suite.T(), err)
}

func (suite *GeneratorTestSuite) TestSingleCatalog() {
	tpl := CatalogTemplate{
		ID:    "abcd",
		Type:  Info,
		Title: "test annon",
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	now := time.Date(2016, time.April, 10, 15, 15, 0, 0, config.Timezone)
	annons, err := gen.FromTemplate(&tpl, &start, nil, nil, nil, &now, nil, 10)
	assert.Len(suite.T(), annons.Time, 1)
	assert.Equal(suite.T(), start, annons.Time[0])
	assert.Nil(suite.T(), err)
	// in the past
	now = now.Add(time.Minute * 10)
	annons, err = gen.FromTemplate(&tpl, &start, nil, nil, nil, &now, nil, 10)
	assert.Nil(suite.T(), annons)
	assert.NoError(suite.T(), err)
}

func (suite *GeneratorTestSuite) TestIsDay() {
	day := time.Date(2016, time.April, 10, 15, 15, 0, 0, config.Timezone)
	assert.True(suite.T(), isDay(&day))
	day = time.Date(2016, time.April, 10, 22, 0, 0, 0, config.Timezone)
	assert.False(suite.T(), isDay(&day))
	day = time.Date(2016, time.April, 10, 22, 1, 0, 0, config.Timezone)
	assert.False(suite.T(), isDay(&day))
	day = time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone)
	assert.False(suite.T(), isDay(&day))
	day = time.Date(2016, time.April, 10, 3, 0, 0, 0, config.Timezone)
	assert.False(suite.T(), isDay(&day))
	day = time.Date(2016, time.April, 10, 7, 0, 0, 0, config.Timezone)
	assert.False(suite.T(), isDay(&day))
}

func (suite *GeneratorTestSuite) TestCatalogRecurrence() {
	tpl := CatalogTemplate{
		ID:    "abcd",
		Type:  Info,
		Title: "test annon",
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 17, 21, 0, 0, config.Timezone)
	now := time.Date(2016, time.April, 10, 15, 15, 0, 0, config.Timezone)
	dayInterval := time.Duration(60) * time.Minute
	nightInterval := time.Duration(90) * time.Minute
	annons, err := gen.FromTemplate(&tpl, &start, &end, &dayInterval, &nightInterval, &now, TemplateParams{}, 10)
	assert.Len(suite.T(), annons.Time, 3)
	assert.Equal(suite.T(), start, annons.Time[0])
	assert.Nil(suite.T(), err)
}

func (suite *GeneratorTestSuite) TestCatalogNightRecurrence() {
	tpl := CatalogTemplate{
		ID:    "abcd",
		Type:  Info,
		Title: "test annon",
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	start := time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 11, 1, 45, 0, 0, config.Timezone)
	now := time.Date(2016, time.April, 10, 15, 15, 0, 0, config.Timezone)
	dayInterval := time.Duration(60) * time.Minute
	nightInterval := time.Duration(90) * time.Minute
	annons, err := gen.FromTemplate(&tpl, &start, &end, &dayInterval, &nightInterval, &now, TemplateParams{}, 10)
	assert.Len(suite.T(), annons.Time, 4)
	assert.Equal(suite.T(), start, annons.Time[0])
	assert.Equal(suite.T(), time.Date(2016, time.April, 11, 0, 0, 0, 0, config.Timezone), annons.Time[2])
	assert.Nil(suite.T(), err)
}

func (suite *GeneratorTestSuite) TestLanguages() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	gen.cache = make(map[string]*CatalogTemplate)
	now := time.Now()
	newTrain := &train.Train{
		Settings: &train.Settings{
			Lang: &map[string]train.LanguageSettings{
				"en": train.LanguageSettings{
					Enabled: true,
				}},
		},
	}
	producer1 := &ProducerMock{}
	producer1.On("Name").Return(Arrival)
	producer1.On("Required", newTrain, mock.AnythingOfType("*train.Train"), mock.AnythingOfType("*time.Time")).Return(true, true, nil)
	event := time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)
	producer1.On("GetTime", newTrain, mock.AnythingOfType("*time.Time")).Return([]time.Time{event}, event, event)
	params := &TemplateParams{}
	producer1.On("BuildParams", newTrain).Return(params, params, true, nil)
	prod1 := Producer(producer1)
	gen.producers[Arrival] = &prod1
	c := &CatalogMock{}
	tpl := CatalogTemplate{
		ID:    "abcd",
		Type:  Arrival,
		Title: "test annon",
		Translations: &map[string]string{
			"en": "tpl_en",
		},
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	en := CatalogTemplate{
		ID:    "tpl_en",
		Type:  Arrival,
		Title: "english annon",
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	c.On("Get", Arrival).Return(&tpl, nil)
	c.On("Get", "tpl_en").Return(&en, nil)
	cat := Catalog(c)
	currentAnnons := &[]*Announcement{}
	res, del, _ := gen.ForTrain(newTrain, nil, currentAnnons, &now, &cat)
	assert.Len(suite.T(), res, 2, "Should contain regular and translated announcement")
	assert.Len(suite.T(), del, 0, "There should be nothing to delete")
}

func (suite *GeneratorTestSuite) TestLangDelete() {
	gen := new(gen)
	gen.producers = make(map[string]*Producer)
	gen.cache = make(map[string]*CatalogTemplate)
	now := time.Now()
	newTrain := &train.Train{
		ID: "1234_test",
		Settings: &train.Settings{
			Lang: &map[string]train.LanguageSettings{
				"en": train.LanguageSettings{
					Enabled: false,
				}},
		},
	}
	producer1 := &ProducerMock{}
	producer1.On("Name").Return(Arrival)
	producer1.On("Required", newTrain, mock.AnythingOfType("*train.Train"), mock.AnythingOfType("*time.Time")).Return(true, true, nil)
	event := time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)
	producer1.On("GetTime", newTrain, mock.AnythingOfType("*time.Time")).Return([]time.Time{event}, event, event)
	params := &TemplateParams{}
	producer1.On("BuildParams", newTrain).Return(params, params, true, nil)
	prod1 := Producer(producer1)
	gen.producers[Arrival] = &prod1
	c := &CatalogMock{}
	tpl := CatalogTemplate{
		ID:    "abcd",
		Type:  Arrival,
		Title: "test annon",
		Translations: &map[string]string{
			"en": "tpl_en",
		},
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	en := CatalogTemplate{
		ID:    "tpl_en",
		Type:  Arrival,
		Title: "english annon",
		Templates: Templates{
			Human: "human text",
			Tts:   "tts text",
			HTML:  "html text",
		},
	}
	c.On("Get", Arrival).Return(&tpl, nil)
	c.On("Get", "tpl_en").Return(&en, nil)
	cat := Catalog(c)
	currentAnnons := &[]*Announcement{&Announcement{ID: "1234_test_arrival_en", Lang: "en"}}
	res, del, _ := gen.ForTrain(newTrain, nil, currentAnnons, &now, &cat)
	assert.Len(suite.T(), res, 1, "Should contain regular and translated announcement")
	assert.Len(suite.T(), del, 1, "There should be old english template to delete")
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
