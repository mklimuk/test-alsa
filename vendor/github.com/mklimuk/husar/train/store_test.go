package train

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/db"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const trainID string = "abcdef1"

type TrainStoreTestSuite struct {
	suite.Suite
	Store   rethinkStore
	Session *r.Session
	skip    bool
	cleanup bool
	db      db.Database
}

func (suite *TrainStoreTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	session, _ := db.Connect("database:28015", 0)
	var d db.Database = "test_train"
	if session != nil && session.IsConnected() {
		log.Debug("Connection established the train will run.")
		db.InitDB(d)
		db.InitTimetable(d)
		db.InitRealtime(d)
		db.InitTrainSettings(d)
		suite.Session = session
		suite.db = d
		suite.Store = newStore(session, suite.db, db.TrainTable, db.RealtimeTable, db.TrainSettingsTable, 5)
		suite.skip = false
		suite.cleanup = true
	} else {
		suite.skip = true
	}
}

func (suite *TrainStoreTestSuite) TearDownSuite() {
	r.DBDrop(suite.db).RunWrite(suite.Session)
}

func (suite *TrainStoreTestSuite) SetupTest() {
	if suite.cleanup {
		CleanupTable(suite.Session, suite.db, db.TrainTable)
		CleanupTable(suite.Session, suite.db, db.RealtimeTable)
		CleanupTable(suite.Session, suite.db, db.TrainSettingsTable)
	}
}

func (suite *TrainStoreTestSuite) TestInitChange() {
	if suite.skip {
		return
	}
	assert.NotPanics(suite.T(), func() { suite.Store.initSettingsChangeCursor() }, "Save all should not panic.")
	assert.NotPanics(suite.T(), func() { suite.Store.initRealtimeChangeCursor() }, "Save all should not panic.")
}

func (suite *TrainStoreTestSuite) TestSaveSettings() {
	if suite.skip {
		return
	}
	ti1 := time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)
	ti2 := time.Date(2016, time.April, 10, 21, 45, 0, 0, config.Timezone)
	ti3 := time.Date(2016, time.April, 10, 21, 59, 0, 0, config.Timezone)
	ti4 := time.Date(2016, time.April, 10, 21, 59, 0, 0, config.Timezone)
	set := []Settings{
		Settings{ID: "t1", FirstEvent: &ti1, Type: TrainInstance},
		Settings{ID: "t1", FirstEvent: &ti2, Type: TrainNumber},
		Settings{ID: "t1", FirstEvent: &ti3, Type: TrainOrder},
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAllSettings(&set) })
	assert.Equal(suite.T(), 1, Count(suite.Session, suite.db, db.TrainSettingsTable), "Inserting elements with the same ID should replace them")
	set[1].ID = "t2"
	set[2].ID = "t3"
	err := suite.Store.SaveAllSettings(&set)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, Count(suite.Session, suite.db, db.TrainSettingsTable), "Inserting elements with different IDs should work just fine")

	last := Settings{ID: "t4", FirstEvent: &ti4, Type: TrainOrder}
	_, err = suite.Store.SaveSettings(&last)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 4, Count(suite.Session, suite.db, db.TrainSettingsTable))
}

func (suite *TrainStoreTestSuite) TestGetSettings() {
	if suite.skip {
		return
	}
	ti1 := time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)
	ti2 := time.Date(2016, time.April, 10, 21, 45, 0, 0, config.Timezone)
	ti3 := time.Date(2016, time.April, 10, 21, 59, 0, 0, config.Timezone)
	ti4 := time.Date(2016, time.October, 29, 0, 0, 0, 0, config.Timezone)
	set := []Settings{
		Settings{ID: "t1", FirstEvent: &ti1, Type: TrainInstance},
		Settings{ID: "t2", FirstEvent: &ti2, Type: TrainNumber},
		Settings{ID: "t3", FirstEvent: &ti3, Type: TrainOrder},
		Settings{ID: "t4", FirstEvent: &ti4, Type: TrainInstance},
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAllSettings(&set) })
	single, err := suite.Store.GetSettings("t4")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), TrainInstance, single.Type)
	ids := []string{"t2", "t1", "t4"}
	multi, err := suite.Store.GetAllSettings(ids...)
	assert.NoError(suite.T(), err)
	// they won't be in order but this is not so important
	assert.Len(suite.T(), *multi, 3)
	single, err = suite.Store.GetSettings("nothere")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), single)
	multi, err = suite.Store.GetAllSettings("nothere", "neither")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), *multi, 0)
}

func (suite *TrainStoreTestSuite) TestSaveRealtime() {
	if suite.skip {
		return
	}
	real := []Realtime{
		Realtime{ID: trainID, Delay: 10, FirstEvent: time.Date(2016, time.April, 10, 15, 23, 0, 0, config.Timezone)},
		Realtime{ID: trainID, Delay: 0, FirstEvent: time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone)},
		Realtime{ID: trainID, Delay: 0, FirstEvent: time.Date(2016, time.April, 10, 15, 27, 0, 0, config.Timezone)},
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAllRealtime(&real) })
	assert.Equal(suite.T(), 1, Count(suite.Session, suite.db, db.RealtimeTable), "Inserting elements with the same ID should replace them")
	real[1].ID = "abcdef2"
	real[2].ID = "abcdef3"
	err := suite.Store.SaveAllRealtime(&real)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, Count(suite.Session, suite.db, db.RealtimeTable), "Inserting elements with different IDs should work just fine")

	last := Realtime{ID: "abcdef4", Delay: 0, FirstEvent: time.Date(2016, time.April, 10, 21, 30, 0, 0, config.Timezone)}
	_, err = suite.Store.SaveRealtime(&last)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 4, Count(suite.Session, suite.db, db.RealtimeTable))

	//test get
	single, err := suite.Store.GetRealtime("abcdef4")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), single)

	single, err = suite.Store.GetRealtime("nothere")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), single)

}

func (suite *TrainStoreTestSuite) TestGetRealtimeBetween() {
	if suite.skip {
		return
	}
	real := []Realtime{
		Realtime{ID: "r4", Delay: 0, FirstEvent: time.Date(2016, time.April, 10, 15, 45, 0, 0, config.Timezone)},
		Realtime{ID: "r1", Delay: 10, FirstEvent: time.Date(2016, time.April, 10, 15, 26, 0, 0, config.Timezone)},
		Realtime{ID: "r3", Delay: 0, FirstEvent: time.Date(2016, time.April, 10, 15, 29, 0, 0, config.Timezone)},
	}
	r1 := time.Date(2016, time.April, 10, 15, 13, 0, 0, config.Timezone)
	r2 := time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone)
	r3 := time.Date(2016, time.April, 10, 15, 29, 0, 0, config.Timezone)
	trains := []Train{
		Train{ID: "r1", FirstEvent: &r1},
		Train{ID: "r2", FirstEvent: &r2},
		Train{ID: "r3", FirstEvent: &r3},
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAllRealtime(&real) })
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAll(&trains) })
	start := time.Date(2016, time.April, 10, 15, 10, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 15, 30, 0, 0, config.Timezone)
	res, err := suite.Store.GetRealtimeBetween(&start, &end)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 2)
	assert.Equal(suite.T(), "r1", res[0].ID)
	assert.Equal(suite.T(), "r3", res[1].ID)
	res, err = suite.Store.GetBetween(&start, &end)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 3)
	assert.Equal(suite.T(), "r2", res[0].ID)
	assert.Equal(suite.T(), "r1", res[1].ID)
	assert.Equal(suite.T(), "r3", res[2].ID)

	last, err := suite.Store.GetRealtime("r4")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "r4", last.ID)

}

func (suite *TrainStoreTestSuite) TestRealtimeChanges() {
	if suite.skip {
		return
	}
	checkpoints := 0
	go func() {
		next, old, err := suite.Store.NextRealtimeChange()
		assert.NoError(suite.T(), err)
		assert.Nil(suite.T(), old)
		assert.Equal(suite.T(), "t1", next.ID)
		checkpoints++
		_, _, _ = suite.Store.NextRealtimeChange()
		checkpoints++
	}()
	rt := Realtime{ID: "t1", FirstEvent: time.Now().Add(time.Minute * 3)}
	suite.Store.SaveRealtime(&rt)
	rt.ID = "t2"
	rt.FirstEvent = time.Now().Add(time.Minute * 10)
	suite.Store.SaveRealtime(&rt)
	assert.Equal(suite.T(), 2, checkpoints)

}

func (suite *TrainStoreTestSuite) TestSettingsChanges() {
	if suite.skip {
		return
	}
	checkpoints := 0
	go func() {
		next, old, err := suite.Store.NextSettingsChange()
		assert.NoError(suite.T(), err)
		assert.Nil(suite.T(), old)
		assert.Equal(suite.T(), "t1", next.ID)
		checkpoints++
		_, _, _ = suite.Store.NextSettingsChange()
		checkpoints++
	}()
	t := time.Now().Add(time.Minute * 3)
	rt := Settings{ID: "t1", FirstEvent: &t}
	suite.Store.SaveSettings(&rt)
	t = time.Now().Add(time.Minute * 10)
	rt.ID = "t2"
	rt.FirstEvent = &t // outside of the window
	suite.Store.SaveSettings(&rt)
	assert.Equal(suite.T(), 2, checkpoints)

}

func (suite *TrainStoreTestSuite) TestSearch() {
	t1 := Train{ID: "searchTest", Name: "testName", TrainID: "abcd", Route: &Route{StartStation: Station{Name: "Radom"}, EndStation: Station{Name: "Lublin"}}}
	t2 := Train{ID: "searchTest2", Name: "test2", TrainID: "efgh", Route: &Route{StartStation: Station{Name: "Gdynia"}, EndStation: Station{Name: "Radom"}}}
	suite.Store.Save(&t1)
	suite.Store.Save(&t2)
	res, err := suite.Store.Search("name")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 1)
	res, err = suite.Store.Search("rad")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), res, 2)
}

func TestTrainStoreTestSuite(t *testing.T) {
	suite.Run(t, new(TrainStoreTestSuite))
}
