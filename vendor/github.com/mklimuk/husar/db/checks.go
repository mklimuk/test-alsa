package db

import (
	"github.com/mklimuk/husar/util"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
)

func databaseExists(d Database) bool {
	return util.SliceContains(getDatabasesList(), d)
}

func tableExists(d Database, t Table) bool {
	return util.SliceContains(getTablesList(d), t)
}

func indexExists(d Database, t Table, i Index) bool {
	return util.SliceContains(getIndexesList(d, t), i)
}

func indexesExist(d Database, t Table, idx []Index) bool {
	result := true
	for _, i := range idx {
		result = result && indexExists(d, t, i)
		if !result {
			return result
		}
	}
	return result
}

func getTablesList(d Database) []string {
	res, err := r.DB(d).TableList().Run(s)
	defer res.Close()
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "db": d, "error": err}).
			Fatal("Could not get tables list.")

		panic(err)
	}
	var tables []string
	err = res.All(&tables)
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "db": d, "error": err}).
			Fatal("Could not parse tables list.")

		panic(err)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "db.checks", "tables": tables}).
			Debug("Found tables.")
	}
	return tables
}

func getDatabasesList() []string {
	res, err := r.DBList().Run(s)
	defer res.Close()
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "error": err}).
			Fatal("Could not read dababases list.")

		panic(err)
	}
	var dbs []string
	err = res.All(&dbs)
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "error": err}).
			Fatal("Could not parse dababases list.")

		panic(err)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "db.checks", "databases": dbs}).
			Debug("Found databases.")
	}
	return dbs
}

func getIndexesList(d Database, t Table) []string {
	res, err := r.DB(d).Table(t).IndexList().Run(s)
	defer res.Close()
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "table": t, "error": err}).
			Fatal("Could not read indexes.")

		panic(err)
	}
	var indexes []string
	err = res.All(&indexes)
	if err != nil {
		log.WithFields(log.Fields{"logger": "db.checks", "table": t, "error": err}).
			Fatal("Could not parse indexes.")

		panic(err)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "db.checks", "indexes": indexes}).
			Debug("Found indexes.")
	}
	return indexes
}

func checkDatabaseStatus(args ...interface{}) bool {
	if len(args) == 0 {
		log.WithFields(log.Fields{"logger": "db.checks", "method": "checkDatabaseStatus", "args": args}).
			Info("Arguments not provided. Nothing to check.")
		return true
	}
	result := false
	l := len(args)
	var d Database
	if l > 0 {
		log.WithFields(log.Fields{"logger": "db.checks", "database": d}).
			Debug("Checking database.")
		d = args[0].(Database)
		result = result || databaseExists(d)
		if !result {
			log.WithFields(log.Fields{"logger": "db.checks", "database": d}).
				Debug("Database not found. Check negative.")
			return false
		}
	}

	var tables map[Table][]Index
	if l > 1 {
		tables = args[1].(map[Table][]Index)
		for t, ind := range tables {
			log.WithFields(log.Fields{"logger": "db.checks", "table": t}).
				Debug("Checking tables.")
			result = result || tableExists(d, t)
			if !result {
				log.WithFields(log.Fields{"logger": "db.checks", "database": d, "table": t}).
					Debug("Table not found. Check negative.")
				return false
			}
			result = result || indexesExist(d, t, ind)
			if !result {
				log.WithFields(log.Fields{"logger": "db.checks", "database": d, "table": t, "index": ind}).
					Debug("Index not found. Check negative.")
				return false
			}
		}
	}
	return result
}
