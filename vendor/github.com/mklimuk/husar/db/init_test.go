package db

import (
	"testing"

	"github.com/mklimuk/husar/util"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBInitTestSuite struct {
	suite.Suite
	Session *r.Session
	skip    bool
	db      Database
}

func (suite *DBInitTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	session, _ := Connect("database:28015", 0)
	if session != nil && session.IsConnected() {
		log.Debug("Established Rethink connection. The test will run.")
		suite.Session = session
		suite.skip = false
		var TestDB Database = "ensure_test"
		suite.db = TestDB
	} else {
		suite.skip = true
	}
}

func (suite *DBInitTestSuite) TearDownSuite() {
	r.DBDrop(suite.db).RunWrite(suite.Session)
}

func (suite *DBInitTestSuite) TestEnsure() {
	checks := map[Table][]Index{
		TrainTable: []Index{
			TrainTimeIndex,
			TrainStationTimeIndex,
		},
		AnnonTable: []Index{
			AnnonTrainIDIndex,
		},
		RealtimeTable: []Index{
			RealtimeTimeIndex,
			RealtimeStationTimeIndex,
		},
		TrainSettingsTable: []Index{},
	}
	if suite.skip {
		return
	}
	_, _ = r.DBDrop(suite.db).RunWrite(suite.Session)
	assert.False(suite.T(), checkDatabaseStatus(suite.db, checks))
	session, err := Connect("database:28015", 0, suite.db, checks)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.NotPanics(suite.T(), func() { Init(suite.db) }, "Ensuring schema should not panic.")
	assert.NotPanics(suite.T(), func() { Init(suite.db) }, "Ensuring schema should not panic when run twice.")
	session, err = Connect("database:28015", 0, suite.db, checks)
	assert.NotNil(suite.T(), session)
	assert.True(suite.T(), session.IsConnected())
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), checkDatabaseStatus(suite.db, checks))
}

func (suite *DBInitTestSuite) TestContains() {
	list := []string{"announcements", "husar", "realtime", "stationFirstLiveEvent"}
	assert.True(suite.T(), util.SliceContains(list, Husar))
	assert.True(suite.T(), util.SliceContains(list, AnnonTable))
	assert.True(suite.T(), util.SliceContains(list, RealtimeStationTimeIndex))
	assert.False(suite.T(), util.SliceContains(list, TrainTable))
}

func TestDBInitTestSuite(t *testing.T) {
	suite.Run(t, new(DBInitTestSuite))
}
