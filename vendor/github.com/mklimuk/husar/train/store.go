package train

import (
	"time"

	"github.com/mklimuk/husar/db"

	r "github.com/dancannon/gorethink"
)

// Store handles annoucements persistence
type Store interface {
	Get(ID string) (*Train, error)
	GetBetween(start, end *time.Time) (res []Train, err error)
	GetAllSettings(IDs ...string) (*[]Settings, error)
	GetRealtime(ID string) (*Realtime, error)
	GetRealtimeBetween(start, end *time.Time) (res []Train, err error)
	GetBetweenWithRealtime(start, end *time.Time) (res []Train, IDs []string, err error)
	NextRealtimeChange() (*Realtime, *Realtime, error)
	NextSettingsChange() (*Settings, *Settings, error)
	SaveAllRealtime(re *[]Realtime) error
	SaveSettings(s *Settings) (*Settings, error)
	Search(query string) (t []*Train, err error)
}

// NewStore announcements store constructor
func NewStore(session *r.Session, db db.Database, trainTable db.Table, realtimeTable db.Table, settingsTable db.Table, windowSize int) Store {
	store := newStore(session, db, trainTable, realtimeTable, settingsTable, windowSize)
	return Store(&store)
}

func newStore(session *r.Session, db db.Database, trainTable db.Table, realtimeTable db.Table, settingsTable db.Table, windowSize int) rethinkStore {
	return rethinkStore{
		session:       session,
		db:            db,
		trainTable:    trainTable,
		realtimeTable: realtimeTable,
		settingsTable: settingsTable,
		windowSize:    windowSize,
	}
}

type rethinkStore struct {
	session        *r.Session
	settingsChange *r.Cursor
	realtimeChange *r.Cursor
	db             db.Database
	trainTable     db.Table
	settingsTable  db.Table
	realtimeTable  db.Table
	windowSize     int
}

//CleanupTable erases table content
func CleanupTable(session *r.Session, d db.Database, table db.Table) {
	r.DB(d).Table(table).Delete().RunWrite(session)
}

//Count counts element in a table
func Count(session *r.Session, d db.Database, table db.Table) int {
	cur, _ := r.DB(d).Table(table).Count().Run(session)
	defer cur.Close()
	var c int
	cur.One(&c)
	return c
}
