package db

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	r "github.com/dancannon/gorethink"
)

var s *r.Session
var address string

// Session returns current RethinkDB r.Session
func Session() (*r.Session, error) {
	if s == nil || !s.IsConnected() {
		if address == "" {
			return nil, errors.New("Session was never properly initialized. Please call Connect() first.")
		}
		if _, err := Connect(address, 10); err != nil {
			log.WithFields(log.Fields{
				"logger": "db.session",
				"error":  err,
			}).Fatal("Could not obtain database connection.")
			return nil, err
		}
		return s, nil
	}
	return s, nil
}

// Connect waits for the database instance to be available and returns a valid r.Session
func Connect(uri string, retries int, args ...interface{}) (*r.Session, error) {
	address = uri
	rethinkConn := false

	log.WithFields(log.Fields{"logger": "db.session", "url": address}).Info("Connecting to RethinkDB.")
	for retries >= 0 {
		if !rethinkConn {
			if rethinkConn = connectToRethink(); !rethinkConn {
				log.WithFields(log.Fields{"logger": "db.session", "connected": rethinkConn}).Info("Connection status.")
				if !rethinkConn {
					log.WithFields(log.Fields{"logger": "db.session"}).Info("Connection not established. Decreasing retries counter.")
					if retries--; retries >= 0 {
						log.Info("Waiting for 10s...")
						time.Sleep(time.Second * 10)
					}
					continue
				}
			}
		}
		if args == nil || len(args) == 0 {
			// no checks specified
			return s, nil
		}
		log.WithFields(log.Fields{"logger": "db.session"}).Info("Verifying database status.")
		// arguments are provided so we need to check database status before returning the session
		if ok := checkDatabaseStatus(args...); ok {
			log.WithFields(log.Fields{"logger": "db.session"}).Info("Database ok. Returning session.")
			return s, nil
		}
		log.WithFields(log.Fields{"logger": "db.session"}).Info("Database not ok. Decreasing retries counter.")
		if retries--; retries >= 0 {
			log.Info("Waiting for 10s...")
			time.Sleep(time.Second * 10)
		}
	}
	return nil, errors.New("Retry counter reached 0. Could not connect to the database.")
}

func connectToRethink() bool {
	var err error
	log.WithFields(log.Fields{"logger": "db.session", "url": address}).Info("Attempting connection to Rethinhdb.")
	if s, err = r.Connect(r.ConnectOpts{
		Address: address,
		MaxIdle: 1,
		MaxOpen: 3,
	}); err != nil {
		log.WithFields(log.Fields{"logger": "db.session", "error": err}).Info("Could not connect to the database.")
		return false
	}
	return s.IsConnected()
}
