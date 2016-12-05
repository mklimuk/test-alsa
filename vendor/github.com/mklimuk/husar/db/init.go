package db

import (
	log "github.com/Sirupsen/logrus"

	r "github.com/dancannon/gorethink"
)

func CreateDatabaseIfNeeded(d Database) {
	// check if the database is present
	if !databaseExists(d) {
		log.WithFields(log.Fields{"logger": "db.init", "db": d}).
			Info("Creating database.")
		// create the database if needed
		_, err := r.DBCreate(d).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{"logger": "db.init", "db": d, "error": err}).
				Fatal("Could not create the database.")
			panic(err)
		}
	}
}

func CreateTableIfNeeded(d Database, table Table) {
	if !tableExists(d, table) {
		log.WithFields(log.Fields{
			"logger": "db.init",
			"table":  table,
		}).Info("Table not found. Creating.")
		_, err := r.DB(d).TableCreate(table).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{
				"logger": "db.init",
				"table":  table,
				"error":  err,
			}).Fatal("Could not create table.")
			panic(err)
		}
	}
}

// Init checks database contents and creates necessary tables and indexes when needed
func Init(d Database) {
	log.WithFields(log.Fields{"logger": "db.init"}).Info("Performing database checks.")
	InitDB(d)
	InitTimetable(d)
	InitAnnouncements(d)
	InitRealtime(d)
	InitTrainSettings(d)
	InitStations(d)
}

//InitDB creates the database if it doesn't exist
func InitDB(d Database) {
	CreateDatabaseIfNeeded(d)
}

//InitTimetable initializes database structures related to timetable storage
func InitTimetable(d Database) {
	log.WithFields(log.Fields{
		"logger": "db.init",
	}).Info("Checking and initializing timetable database structures.")
	CreateTableIfNeeded(d, TrainTable)
	updateTimetableIndexes(d)
}

//InitRealtime initializes database structures related to realtime timetable storage
func InitRealtime(d Database) {
	log.WithFields(log.Fields{
		"logger": "db.init",
	}).Info("Checking and initializing realtime database structures.")
	CreateTableIfNeeded(d, RealtimeTable)
	updateRealtimeIndexes(d)
}

//InitAnnouncements initializes database structures related to announcements storage
func InitAnnouncements(d Database) {
	log.WithFields(log.Fields{
		"logger": "db.init",
	}).Info("Checking and initializing announcements database structures.")
	CreateTableIfNeeded(d, AnnonTable)
	updateAnnonIndexes(d)
}

//InitTrainSettings initializes database structures related to train settings storage
func InitTrainSettings(d Database) {
	log.WithFields(log.Fields{
		"logger": "db.init",
	}).Info("Checking and initializing train settings structures.")
	CreateTableIfNeeded(d, TrainSettingsTable)
	updateSettingsIndexes(d)
}

//InitStations initializes database structures related to stations settings storage
func InitStations(d Database) {
	log.WithFields(log.Fields{
		"logger": "db.init",
	}).Info("Checking and initializing stations settings structures.")
	CreateTableIfNeeded(d, StationsTable)
	updateSettingsIndexes(d)
}

func updateAnnonIndexes(d Database) {
	if !indexExists(d, AnnonTable, AnnonTrainIDIndex) {

		log.WithFields(log.Fields{"logger": "db.init", "table": AnnonTable, "index": AnnonTrainIDIndex}).
			Info("Index not found. Creating.")

		_, err := r.DB(d).Table(AnnonTable).IndexCreate(AnnonTrainIDIndex).RunWrite(s)
		if err != nil {

			log.WithFields(log.Fields{"logger": "db.init", "table": AnnonTable, "index": AnnonTrainIDIndex, "error": err}).
				Fatal("Could not create index.")

			panic(err)
		}
	}
	if !indexExists(d, AnnonTable, AnnonTimeIndex) {

		log.WithFields(log.Fields{"logger": "db.init", "table": AnnonTable, "index": AnnonTrainIDIndex}).
			Info("Index not found. Creating.")

		_, err := r.DB(d).Table(AnnonTable).IndexCreate(AnnonTimeIndex).RunWrite(s)
		if err != nil {

			log.WithFields(log.Fields{"logger": "db.init", "table": AnnonTable, "index": AnnonTimeIndex, "error": err}).
				Fatal("Could not create index.")

			panic(err)
		}
	}
}

func updateSettingsIndexes(d Database) {

}

func updateTimetableIndexes(d Database) {
	if !indexExists(d, TrainTable, TrainTimeIndex) {
		log.WithFields(log.Fields{
			"logger": "db.init",
			"table":  TrainTable,
			"index":  TrainTimeIndex,
		}).Info("Index not found. Creating.")
		_, err := r.DB(d).Table(TrainTable).IndexCreate(TrainTimeIndex).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{
				"logger": "db.init",
				"table":  TrainTable,
				"index":  TrainTimeIndex,
				"error":  err,
			}).Fatal("Could not create index.")
			panic(err)
		}
	}
	if !indexExists(d, TrainTable, TrainStationTimeIndex) {
		log.WithFields(log.Fields{
			"logger": "db.init",
			"table":  TrainTable,
			"index":  TrainStationTimeIndex,
		}).Info("Index not found. Creating.")
		_, err := r.DB(d).Table(TrainTable).IndexCreateFunc(TrainStationTimeIndex, func(row r.Term) interface{} {
			return []interface{}{row.Field(StationIDField), row.Field(FirstEventField)}
		}).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{
				"logger": "db.init",
				"table":  TrainTable,
				"index":  TrainTimeIndex,
				"error":  err,
			}).Fatal("Could not create index.")
			panic(err)
		}
	}
}

func updateRealtimeIndexes(d Database) {
	if !indexExists(d, RealtimeTable, RealtimeTimeIndex) {
		log.WithFields(log.Fields{
			"logger": "db.init",
			"table":  RealtimeTable,
			"index":  RealtimeTimeIndex,
		}).Info("Index not found. Creating.")
		_, err := r.DB(d).Table(RealtimeTable).IndexCreate(RealtimeTimeIndex).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{
				"logger": "db.init",
				"table":  RealtimeTable,
				"index":  RealtimeTimeIndex,
				"error":  err,
			}).Fatal("Could not create index.")
			panic(err)
		}
	}
	if !indexExists(d, RealtimeTable, RealtimeStationTimeIndex) {
		log.WithFields(log.Fields{
			"logger": "db.init",
			"table":  RealtimeTable,
			"index":  RealtimeStationTimeIndex,
		}).Info("Index not found. Creating.")
		_, err := r.DB(d).Table(RealtimeTable).IndexCreateFunc(RealtimeStationTimeIndex, func(row r.Term) interface{} {
			return []interface{}{row.Field(StationIDField), row.Field(RealtimeFirstEventField)}
		}).RunWrite(s)
		if err != nil {
			log.WithFields(log.Fields{
				"logger": "db.init",
				"table":  RealtimeTable,
				"index":  RealtimeStationTimeIndex,
				"error":  err,
			}).Fatal("Could not create index.")
			panic(err)
		}
	}
}
