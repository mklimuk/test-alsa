package db

// Database represents the database
type Database string

// Index represents table indexes
type Index string

// Table represents tables
type Table string

// Field represents fields for complex indexes
type Field string

// RethinkDB related constants
const (
	Husar Database = "husar"
	// tables
	AnnonTable         Table = "announcements"
	TrainTable         Table = "timetable"
	RealtimeTable      Table = "realtime"
	TrainSettingsTable Table = "train_settings"
	AnnonSettingsTable Table = "annon_settings"

	StationsTable Table = "stations"
	// indexes
	AnnonTrainIDIndex        Index = "trainId"
	AnnonTimeIndex           Index = "first"
	RealtimeTimeIndex        Index = "firstLiveEvent"
	SettingsTimeIndex        Index = "firstEvent"
	TrainTimeIndex           Index = "firstEvent"
	RealtimeStationTimeIndex Index = "stationFirstLiveEvent"
	TrainStationTimeIndex    Index = "stationFirstEvent"
	// fields
	StationIDField          Field = "stationId"
	FirstEventField         Field = "firstEvent"
	RealtimeFirstEventField Field = "firstLiveEvent"
)
