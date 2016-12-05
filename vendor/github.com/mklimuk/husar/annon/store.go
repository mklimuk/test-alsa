package annon

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/mklimuk/husar/db"

	r "github.com/dancannon/gorethink"
)

// Store handles annoucements persistence
type Store interface {
	Get(id string) (*Announcement, error)
	GetForTrain(trainID string) (*[]*Announcement, error)
	GetBetween(start, end *time.Time) (res []Announcement, err error)
	Save(a *Announcement) (*Announcement, error)
	Delete(id string) error
	DeleteForTrain(trainID string) error
	SaveAll(a *[]*Announcement) error
	NextChange() (*Announcement, *Announcement, error)
}

// NewStore announcements store constructor
func NewStore(session *r.Session, db db.Database, table db.Table, windowSize int) Store {
	store := rethinkStore{
		session:    session,
		db:         db,
		table:      table,
		windowSize: windowSize,
	}
	return Store(&store)
}

type rethinkStore struct {
	session    *r.Session
	change     *r.Cursor
	db         db.Database
	table      db.Table
	windowSize int
}

func (store *rethinkStore) Get(id string) (*Announcement, error) {
	resp, err := r.DB(store.db).Table(store.table).Get(id).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "Get", "template": id}).
			WithError(err).Error("Error loading template.")
		return nil, err
	}
	a := new(Announcement)
	err = resp.One(a)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annon.store", "method": "Get", "template": id}).
			WithError(err).Error("Error parsing template.")
		return nil, err
	}
	return a, err
}

func (store *rethinkStore) Save(a *Announcement) (*Announcement, error) {
	resp, err := r.DB(store.db).Table(store.table).Insert(a, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "Save", "announcement": a}).
			WithError(err).Error("Error occured while saving announcement.")
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "Save", "announcement": a, "result": resp}).
			Debug("Database save result.")
	}
	return a, err
}

func (store *rethinkStore) Delete(ID string) error {
	_, err := r.DB(store.db).Table(store.table).Get(ID).Delete().RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "Delete", "announcement": ID}).
			WithError(err).Error("Error occured while deleting announcement.")
	}
	return err
}

func (store *rethinkStore) SaveAll(a *[]*Announcement) error {
	res, err := r.DB(store.db).Table(store.table).Insert(*a, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "SaveAll", "announcement": a}).
			WithError(err).Error("Error occured while saving announcements.")
		return err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "SavelAll", "announcements": a, "result": res}).
			Debug("Database save result.")
	}
	return err
}

func (store *rethinkStore) DeleteForTrain(trainID string) error {
	res, err := r.DB(store.db).Table(store.table).GetAllByIndex(string(db.AnnonTrainIDIndex), trainID).Delete().RunWrite(store.session)
	log.WithFields(log.Fields{
		"logger": "annon.store", "method": "DeleteForTrain", "train": trainID, "result": res}).
		Debug("Database delete result.")
	return err
}

func (store *rethinkStore) GetForTrain(trainID string) (*[]*Announcement, error) {
	resp, err := r.DB(store.db).Table(store.table).GetAllByIndex(string(db.AnnonTrainIDIndex), trainID).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annon.store", "method": "GetForTrain", "train": trainID}).
			WithError(err).Error("Error loading announcements for train.")
		return nil, err
	}
	a := new([]*Announcement)
	err = resp.All(a)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annon.store", "method": "GetForTrain", "train": trainID}).
			WithError(err).Error("Error parsing announcements for train.")
		return nil, err
	}
	return a, err
}

func (store *rethinkStore) GetBetween(start, end *time.Time) (annons []Announcement, err error) {
	var res *r.Cursor
	if res, err = r.DB(store.db).Table(store.table).
		Between(start, end, r.BetweenOpts{Index: db.AnnonTimeIndex}).
		Run(store.session); err != nil {
		log.WithFields(log.Fields{"logger": "model.annon.store", "method": "GetBetween", "start": *start, "end": end}).
			WithError(err).Error("Error executing query on the announcements table.")
		return
	}
	if err = res.All(&annons); err != nil {
		log.WithFields(log.Fields{"logger": "model.annon.store", "method": "GetBetween", "start": *start, "end": end}).
			WithError(err).Error("Error parsing results from the announcements table.")
		return
	}
	return
}

func (store *rethinkStore) initChangeCursor() error {
	log.WithFields(log.Fields{"logger": "model.annon.store", "method": "initChangeCursor", "windowSize": store.windowSize}).
		Info("Initializing change cursor for announcements timetable changes.")
	res, err := r.DB(store.db).Table(store.table).Changes(r.ChangesOpts{IncludeInitial: false}).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "model.annon.store", "method": "initChangeCursor"}).
			WithError(err).Error("Error while initializing announcements table changes cursor.")
		return err
	}
	store.change = res
	return nil
}

func (store *rethinkStore) NextChange() (*Announcement, *Announcement, error) {
	if store.change == nil {
		err := store.initChangeCursor()
		if err != nil {
			return nil, nil, err
		}
	}
	var next bool
	changes := new(Changes)
	if next = store.change.Next(changes); !next {
		err := store.change.Err()
		log.WithFields(log.Fields{"logger": "train.realtime", "method": "NextChange"}).
			Error("Error while reading from realtime table changes cursor.")
		return nil, nil, err
	}
	return changes.NewVal, changes.OldVal, nil
}
