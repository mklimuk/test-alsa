package train

import (
	"strconv"
	"testing"
	"time"

	"golang.org/x/text/language"

	"github.com/mklimuk/husar/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TrainUtilsTestSuite struct {
	suite.Suite
}

func (suite *TrainUtilsTestSuite) TestServiceChecker() {
	train := Train{}
	test, _ := HasService(&train, REZO)
	assert.Equal(suite.T(), false, test)
	train.Services = &[]Service{
		Service{ID: "REZO", Name: "rezerwacja"},
		Service{ID: "KOND", Name: "konduktorskie", Carriage: []string{"10"}},
	}
	test, _ = HasService(&train, REZO)
	assert.Equal(suite.T(), true, test)
	test, car := HasService(&train, KOND)
	assert.Equal(suite.T(), true, test)
	assert.Equal(suite.T(), []string{"10"}, *car)
}

func (suite *TrainUtilsTestSuite) TestOptionsChecker() {
	assert.NotPanics(suite.T(), func() {
		HasAnnonOption(nil, ArrivalText)
	})
	train := Train{}
	test, _ := HasAnnonOption(&train, ArrivalText)
	assert.Equal(suite.T(), false, test)
	train.Settings = &Settings{
		Audio: &map[AnnonOption]string{ArrivalText: "Nietypowy tekst", Pause: "false"},
	}
	test, value := HasAnnonOption(&train, ArrivalText)
	assert.Equal(suite.T(), true, test)
	assert.Equal(suite.T(), "Nietypowy tekst", value)
	test, value = HasAnnonOption(&train, Pause)
	assert.Equal(suite.T(), true, test)
	v, _ := strconv.ParseBool(value)
	assert.Equal(suite.T(), false, v)
}

func (suite *TrainUtilsTestSuite) TestMiddleExtraction() {
	r := routes["middle"]
	st, through := FollowingStations(r, true)
	assert.Equal(suite.T(), "Station6, Station7", through)
	assert.Equal(suite.T(), 6, st[0].ID)
}

func (suite *TrainUtilsTestSuite) TestEndExtraction() {
	r := routes["end"]
	st, through := FollowingStations(r, true)
	assert.Equal(suite.T(), "", through)
	assert.Equal(suite.T(), 0, len(st))
}

func (suite *TrainUtilsTestSuite) TestStartExtraction() {
	r := routes["start"]
	st, through := FollowingStations(r, true)
	assert.Equal(suite.T(), "Station4, Station5", through)
	assert.Equal(suite.T(), 2, len(st))
}

func (suite *TrainUtilsTestSuite) TestEmptyExtraction() {
	r := routes["empty"]
	st, through := FollowingStations(r, true)
	assert.Equal(suite.T(), through, "")
	assert.Equal(suite.T(), len(st), 0)
}

func (suite *TrainUtilsTestSuite) TestSettingsMap() {
	a := Settings{ID: "a"}
	b := Settings{ID: "b"}
	s := []Settings{a, b}
	set := BuildSettingsMap(&s)
	assert.Len(suite.T(), set, 2)
	assert.Equal(suite.T(), a, set["a"])
	assert.Equal(suite.T(), b, set["b"])
}

func (suite *TrainUtilsTestSuite) TestCompareSettings() {
	a := Train{}
	b := Train{}
	assert.True(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	a.Settings = &Settings{}
	assert.False(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	a.Settings.Audio = &map[AnnonOption]string{}
	assert.False(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	b.Settings = &Settings{}
	assert.False(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	b.Settings.Audio = &map[AnnonOption]string{}
	assert.True(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	(*a.Settings.Audio)[ArrivalText] = "abc"
	assert.False(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	(*b.Settings.Audio)[ArrivalText] = "abc"
	assert.True(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))
	(*b.Settings.Audio)[ArrivalText] = "abcd"
	assert.False(suite.T(), CompareAudioSettings(&a, &b, ArrivalText))

}

func (suite *TrainUtilsTestSuite) TestCopy() {
	copy := new(Train)
	ev := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	orig := Train{
		ID:      "abcdef12345",
		TrainID: "Train1",
		Delay:   0,
		Arrival: &TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
			Platform: "II",
			Track:    1,
		},
		Departure: &TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone),
			Platform: "II",
			Track:    1,
		},
		FirstEvent: &ev,
		Route: &Route{
			CurrentStationPositionOnRoute: 5,
			StartStation: Station{
				ID:            2,
				Name:          "Station2",
				Important:     true,
				NumberOnRoute: 1,
			},
			EndStation: Station{
				ID:            3,
				Name:          "Station3",
				Important:     true,
				NumberOnRoute: 15,
			},
			Stations: []Station{
				Station{
					ID:            2,
					Name:          "Station2",
					Important:     true,
					NumberOnRoute: 1,
				},
				Station{
					ID:            4,
					Name:          "Station4",
					Important:     false,
					NumberOnRoute: 3,
				},
			},
		},
		Carrier:     "IC",
		Category:    "TLK",
		Day:         time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone),
		OrderID:     100,
		StationID:   1,
		StationName: "Station1",
		Name:        "LUNA",
	}
	copy = Copy(&orig, copy)
	copy.Arrival.Platform = "test"
	copy.Departure.Platform = "test"
	copy.Route.Stations[0].Name = "test"
	copy.Route.StartStation.Name = "test"
	copy.Category = "cat"
	assert.NotEqual(suite.T(), orig.Category, copy.Category)
	assert.NotEqual(suite.T(), orig.Arrival.Platform, copy.Arrival.Platform)
	assert.NotEqual(suite.T(), orig.Departure.Platform, copy.Departure.Platform)
	assert.Len(suite.T(), copy.Route.Stations, 2)
	assert.NotEqual(suite.T(), orig.Route.Stations[0].Name, copy.Route.Stations[0].Name)
	assert.NotEqual(suite.T(), orig.Route.StartStation.Name, copy.Route.StartStation.Name)
}

func (suite *TrainUtilsTestSuite) TestGetCategory() {
	t := Train{Category: "abcd"}
	assert.Equal(suite.T(), "abcd", GetCategory(&t))
	t.Settings = &Settings{}
	assert.Equal(suite.T(), "abcd", GetCategory(&t))
	t.Settings.Overrides = &map[string]string{}
	assert.Equal(suite.T(), "abcd", GetCategory(&t))
	(*t.Settings.Overrides)[Category] = "dcba"
	assert.Equal(suite.T(), "dcba", GetCategory(&t))
}

func (suite *TrainUtilsTestSuite) TestCompareI18n() {
	t1 := Train{}
	t2 := Train{}
	assert.True(suite.T(), CompareI18n(&t1, &t2, ArrivalText))
	t1.Settings = &Settings{
		Lang: &map[string]LanguageSettings{
			"en": LanguageSettings{
				Enabled: true,
			},
		},
	}
	t2.Settings = &Settings{
		Lang: &map[string]LanguageSettings{
			"en": LanguageSettings{
				Enabled: false,
			},
		},
	}
	assert.True(suite.T(), CompareI18n(&t1, &t2, ArrivalText))
	t1.Settings = &Settings{
		Lang: &map[string]LanguageSettings{
			"pl": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					ArrivalText: "testujemy",
				},
			},
			"en": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					ArrivalText: "just testing",
				},
			},
		},
	}
	assert.False(suite.T(), CompareI18n(&t1, &t2, ArrivalText))
	t2.Settings = &Settings{
		Lang: &map[string]LanguageSettings{
			"pl": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					ArrivalText: "testujemy",
				},
			},
			"en": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					ArrivalText: "just testing",
				},
			},
		},
	}
	assert.True(suite.T(), CompareI18n(&t1, &t2, ArrivalText))
	options := (*t2.Settings.Lang)["en"].I18n
	(*options)[ArrivalText] = "new one"
	assert.False(suite.T(), CompareI18n(&t1, &t2, ArrivalText))
}

func (suite *TrainUtilsTestSuite) TestHasI18n() {
	t1 := Train{}
	test, _ := HasI18nOption(&t1, ArrivalText, language.Polish)
	assert.False(suite.T(), test)
	t1.Settings = &Settings{
		Lang: &map[string]LanguageSettings{
			"pl": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					ArrivalText: "testujemy",
				},
			},
			"en": LanguageSettings{
				Enabled: true,
				I18n: &map[AnnonOption]string{
					DepartureText: "just testing",
				},
			},
		},
	}
	test, val := HasI18n(&t1)
	assert.True(suite.T(), test)
	assert.Len(suite.T(), *val, 2)
	test, text := HasI18nOption(&t1, ArrivalText, language.Polish)
	assert.True(suite.T(), test)
	assert.Equal(suite.T(), "testujemy", text)
	test, text = HasI18nOption(&t1, ArrivalText, language.English)
	assert.False(suite.T(), test)
	test, text = HasI18nOption(&t1, DepartureText, language.English)
	assert.True(suite.T(), test)
	assert.Equal(suite.T(), "just testing", text)
}

func (suite *TrainUtilsTestSuite) TestCompareOverrides() {
	t1 := Train{}
	t2 := Train{Settings: &Settings{Overrides: &map[string]string{}}}
	assert.False(suite.T(), CompareOverrides(&t1, &t2))
	t1.Settings = &Settings{Overrides: &map[string]string{}}
	assert.True(suite.T(), CompareOverrides(&t1, &t2))
	(*t1.Settings.Overrides)["category"] = "special"
	assert.False(suite.T(), CompareOverrides(&t1, &t2))
	(*t2.Settings.Overrides)["category"] = "special"
	assert.True(suite.T(), CompareOverrides(&t1, &t2))
}

func TestTrainUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(TrainUtilsTestSuite))
}

var routes = map[string]*Route{
	"start": &Route{
		CurrentStationPositionOnRoute: 1,
		StartStation: Station{
			ID:            1,
			Name:          "Station1",
			Important:     true,
			NumberOnRoute: 1,
		},
		EndStation: Station{
			ID:            3,
			Name:          "Station3",
			Important:     true,
			NumberOnRoute: 10,
		},
		Stations: []Station{
			Station{
				ID:            1,
				Name:          "Station1",
				Important:     true,
				NumberOnRoute: 1,
			},
			Station{
				ID:            4,
				Name:          "Station4",
				Important:     true,
				NumberOnRoute: 3,
			},
			Station{
				ID:            5,
				Name:          "Station5",
				Important:     true,
				NumberOnRoute: 7,
			},
			Station{
				ID:            6,
				Name:          "Station6",
				Important:     false,
				NumberOnRoute: 9,
			},
			Station{
				ID:            3,
				Name:          "Station3",
				Important:     false,
				NumberOnRoute: 10,
			},
		},
	},
	"middle": &Route{
		CurrentStationPositionOnRoute: 5,
		StartStation: Station{
			ID:            2,
			Name:          "Station2",
			Important:     true,
			NumberOnRoute: 1,
		},
		EndStation: Station{
			ID:            3,
			Name:          "Station3",
			Important:     true,
			NumberOnRoute: 15,
		},
		Stations: []Station{
			Station{
				ID:            2,
				Name:          "Station2",
				Important:     true,
				NumberOnRoute: 1,
			},
			Station{
				ID:            4,
				Name:          "Station4",
				Important:     false,
				NumberOnRoute: 3,
			},
			Station{
				ID:            1,
				Name:          "Station1",
				Important:     true,
				NumberOnRoute: 5,
			},
			Station{
				ID:            5,
				Name:          "Station5",
				Important:     false,
				NumberOnRoute: 7,
			},
			Station{
				ID:            6,
				Name:          "Station6",
				Important:     true,
				NumberOnRoute: 9,
			},
			Station{
				ID:            7,
				Name:          "Station7",
				Important:     true,
				NumberOnRoute: 11,
			},
			Station{
				ID:            3,
				Name:          "Station3",
				Important:     true,
				NumberOnRoute: 15,
			},
		},
	},
	"end": &Route{
		CurrentStationPositionOnRoute: 10,
		StartStation: Station{
			ID:            2,
			Name:          "Station2",
			Important:     true,
			NumberOnRoute: 1,
		},
		EndStation: Station{
			ID:            1,
			Name:          "Station1",
			Important:     true,
			NumberOnRoute: 10,
		},
		Stations: []Station{
			Station{
				ID:            2,
				Name:          "Station2",
				Important:     true,
				NumberOnRoute: 1,
			},
			Station{
				ID:            4,
				Name:          "Station4",
				Important:     false,
				NumberOnRoute: 3,
			},
			Station{
				ID:            7,
				Name:          "Station7",
				Important:     true,
				NumberOnRoute: 8,
			},
			Station{
				ID:            1,
				Name:          "Station1",
				Important:     true,
				NumberOnRoute: 10,
			},
		},
	},
	"empty": &Route{},
}
