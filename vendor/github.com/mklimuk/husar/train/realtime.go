package train

import (
	"fmt"
	"github.com/mklimuk/husar/db"
	"time"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
)

// Realtime stores information about realtime events for a running train
type Realtime struct {
	ID         string    `gorethink:"id" json:"id"`
	StationID  int       `gorethink:"stationId" json:"stationId"`
	FirstEvent time.Time `gorethink:"firstLiveEvent" json:"firstLiveEvent"`
	Delay      int       `gorethink:"delay" json:"delay"`
}

// RealtimeChanges represents changefeeds data from database
type RealtimeChanges struct {
	NewVal *Realtime `gorethink:"new_val,omitempty"`
	OldVal *Realtime `gorethink:"old_val,omitempty"`
}

/*
Store methods related to Realtime
*/

func (store *rethinkStore) GetRealtimeBetween(start, end *time.Time) (trains []Train, err error) {
	var res *r.Cursor
	if res, err = r.DB(store.db).Table(store.realtimeTable).
		Between(start, end, r.BetweenOpts{Index: db.RealtimeTimeIndex}).
		EqJoin("id", r.DB(store.db).Table(store.trainTable)).Zip().
		OrderBy(db.RealtimeTimeIndex).
		Run(store.session); err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "GetRealtimeBetween", "start": start, "end": end}).
			WithError(err).Error("Error executing query on the realtime table.")
		return
	}
	if err = res.All(&trains); err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "GetRealtimeBetween", "start": start, "end": end}).
			WithError(err).Error("Error parsing results from the realtime table.")
		return
	}
	return
}

func (store *rethinkStore) initRealtimeChangeCursor() error {
	log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "initRealtimeChangeCursor", "windowSize": store.windowSize}).
		Info("Initializing change cursor for realtime timetable changes.")
	res, err := r.DB(store.db).Table(store.realtimeTable).
		Changes(r.ChangesOpts{IncludeInitial: false}).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "initRealtimeChangeCursor"}).
			WithError(err).Error("Error while initializing train table changes cursor.")
		return err
	}
	store.realtimeChange = res
	return nil
}

func (store *rethinkStore) SaveRealtime(rt *Realtime) (*Realtime, error) {
	resp, err := r.DB(store.db).Table(store.realtimeTable).Insert(rt, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "SaveRealtime", "realtime": fmt.Sprintf("%+v", rt)}).
			WithError(err).Error("Error occured while saving realtime object.")
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "SaveRealtime", "realtime": fmt.Sprintf("%+v", rt), "result": fmt.Sprintf("%+v", resp)}).
			Debug("Database save result.")
	}
	return rt, err
}

func (store *rethinkStore) SaveAllRealtime(re *[]Realtime) error {
	res, err := r.DB(store.db).Table(store.realtimeTable).Insert(*re, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "SaveAllRealtime", "realtime": re}).
			WithError(err).Error("Error occured while saving realtime objects.")
		return err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "SaveAllRealtime", "realtime": re, "result": res}).
			Debug("Database save result.")
	}
	return err
}

func (store *rethinkStore) GetRealtime(ID string) (*Realtime, error) {
	resp, err := r.DB(store.db).Table(store.realtimeTable).Get(ID).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "GetRealtime", "train": ID}).
			WithError(err).Error("Error while retrieving realtime train.")
		return nil, err
	}
	if resp.IsNil() {
		return nil, err
	}
	a := new(Realtime)
	err = resp.One(a)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "GetRealtime", "train": ID}).
			WithError(err).Error("Error while parsing realtime train result.")
		return nil, err
	}
	return a, err
}

func (store *rethinkStore) NextRealtimeChange() (*Realtime, *Realtime, error) {
	if store.realtimeChange == nil {
		err := store.initRealtimeChangeCursor()
		if err != nil {
			return nil, nil, err
		}
	}
	var next bool
	changes := new(RealtimeChanges)
	if next = store.realtimeChange.Next(changes); !next {
		err := store.realtimeChange.Err()
		log.WithFields(log.Fields{"logger": "train.store.realtime", "method": "NextRealtimeChange"}).
			WithError(err).Error("Error while reading from realtime table changes cursor.")
		return nil, nil, err
	}
	return changes.NewVal, changes.OldVal, nil
}
