package test

import (
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/train"
)

// GetTrain returns a copy of one of predefined trains
func GetTrain(id string) *train.Train {
	orig := trains[id]
	if orig == nil {
		return nil
	}
	c := new(train.Train)
	return train.Copy(orig, c)
}

var passThrough = time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
var starting = time.Date(2016, time.April, 10, 9, 8, 0, 0, config.Timezone)
var ending = time.Date(2016, time.April, 10, 0, 57, 0, 0, config.Timezone)

var trains = map[string]*train.Train{
	"passThrough": &train.Train{
		ID:      "abcdef12345",
		TrainID: "Train1",
		Delay:   0,
		Arrival: &train.TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
			Platform: "II",
			Track:    1,
		},
		Departure: &train.TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone),
			Platform: "II",
			Track:    1,
		},
		FirstEvent: &passThrough,
		Route: &train.Route{
			CurrentStationPositionOnRoute: 5,
			StartStation: train.Station{
				ID:            2,
				Name:          "Station2",
				Important:     true,
				NumberOnRoute: 1,
			},
			EndStation: train.Station{
				ID:            3,
				Name:          "Station3",
				Important:     true,
				NumberOnRoute: 15,
			},
			Stations: []train.Station{
				train.Station{
					ID:            2,
					Name:          "Station2",
					Important:     true,
					NumberOnRoute: 1,
				},
				train.Station{
					ID:            4,
					Name:          "Station4",
					Important:     false,
					NumberOnRoute: 3,
				},
				train.Station{
					ID:            1,
					Name:          "Station1",
					Important:     true,
					NumberOnRoute: 5,
				},
				train.Station{
					ID:            5,
					Name:          "Station5",
					Important:     false,
					NumberOnRoute: 7,
				},
				train.Station{
					ID:            6,
					Name:          "Station6",
					Important:     true,
					NumberOnRoute: 9,
				},
				train.Station{
					ID:            7,
					Name:          "Station7",
					Important:     true,
					NumberOnRoute: 11,
				},
				train.Station{
					ID:            3,
					Name:          "Station3",
					Important:     true,
					NumberOnRoute: 15,
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
	},
	"passThroughNoIntermediaryStations": &train.Train{},
	"starting": &train.Train{
		ID:      "abcdef12345",
		TrainID: "Train2",
		Delay:   0,
		Departure: &train.TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 9, 8, 0, 0, config.Timezone),
			Platform: "II",
			Track:    1,
		},
		FirstEvent: &starting,
		Route: &train.Route{
			CurrentStationPositionOnRoute: 1,
			StartStation: train.Station{
				ID:            1,
				Name:          "Station1",
				Important:     true,
				NumberOnRoute: 1,
			},
			EndStation: train.Station{
				ID:            3,
				Name:          "Station3",
				Important:     true,
				NumberOnRoute: 10,
			},
			Stations: []train.Station{
				train.Station{
					ID:            1,
					Name:          "Station1",
					Important:     true,
					NumberOnRoute: 1,
				},
				train.Station{
					ID:            4,
					Name:          "Station4",
					Important:     true,
					NumberOnRoute: 3,
				},
				train.Station{
					ID:            5,
					Name:          "Station5",
					Important:     true,
					NumberOnRoute: 7,
				},
				train.Station{
					ID:            6,
					Name:          "Station6",
					Important:     false,
					NumberOnRoute: 9,
				},
				train.Station{
					ID:            3,
					Name:          "Station3",
					Important:     false,
					NumberOnRoute: 10,
				},
			},
		},
		Carrier:     "KM",
		Category:    "Os",
		Day:         time.Date(2016, time.April, 10, 0, 0, 0, 0, config.Timezone),
		OrderID:     101,
		StationID:   1,
		StationName: "Station1",
		Name:        "RADOMIAK",
	},
	"startingNoIntermediaryStations": &train.Train{},
	"ending": &train.Train{
		ID:      "abcdef12345",
		TrainID: "Train3",
		Delay:   0,
		Arrival: &train.TimetableEvent{
			Time:     time.Date(2016, time.April, 10, 0, 57, 0, 0, config.Timezone),
			Platform: "III",
			Track:    8,
		},
		FirstEvent: &ending,
		Route: &train.Route{
			CurrentStationPositionOnRoute: 10,
			StartStation: train.Station{
				ID:            2,
				Name:          "Station2",
				Important:     true,
				NumberOnRoute: 1,
			},
			EndStation: train.Station{
				ID:            1,
				Name:          "Station1",
				Important:     true,
				NumberOnRoute: 10,
			},
			Stations: []train.Station{
				train.Station{
					ID:            2,
					Name:          "Station2",
					Important:     true,
					NumberOnRoute: 1,
				},
				train.Station{
					ID:            4,
					Name:          "Station4",
					Important:     false,
					NumberOnRoute: 3,
				},
				train.Station{
					ID:            7,
					Name:          "Station7",
					Important:     true,
					NumberOnRoute: 8,
				},
				train.Station{
					ID:            1,
					Name:          "Station1",
					Important:     true,
					NumberOnRoute: 10,
				},
			},
		},
		Carrier:     "KM",
		Category:    "Os",
		Day:         time.Date(2016, time.April, 9, 0, 0, 0, 0, config.Timezone),
		OrderID:     101,
		StationID:   1,
		StationName: "Station1",
	},
}
