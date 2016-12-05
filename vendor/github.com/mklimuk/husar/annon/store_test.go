package annon

import (
	"testing"
	"time"

	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/db"

	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const trainID string = "abcdef12345"

type AnnonStoreTestSuite struct {
	suite.Suite
	Store   Store
	Session *r.Session
	skip    bool
	d       db.Database
}

func (suite *AnnonStoreTestSuite) SetupSuite() {
	var d db.Database = "annon_test"
	suite.d = d
	session, _ := db.Connect("database:28015", 0)
	if session != nil && session.IsConnected() {
		db.InitDB(d)
		db.InitAnnouncements(d)
		suite.Session = session
		suite.Store = NewStore(session, d, db.AnnonTable, 5)
		suite.skip = false
	} else {
		suite.skip = true
	}
}

func (suite *AnnonStoreTestSuite) SetupTest() {
	cleanupAnnons(suite.Session, suite.d)
}

func (suite *AnnonStoreTestSuite) TestSaveAll() {
	if suite.skip {
		return
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAll(annons) }, "Save all should not panic.")
	assert.Equal(suite.T(), 3, countAnnons(suite.Session, suite.d))
}

func (suite *AnnonStoreTestSuite) TestGetForTrain() {
	if suite.skip {
		return
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAll(annons) }, "Save all should not panic.")
	assert.Equal(suite.T(), 3, countAnnons(suite.Session, suite.d))
	forTrain, _ := suite.Store.GetForTrain(trainID)
	assert.Len(suite.T(), *forTrain, 2)
}

func (suite *AnnonStoreTestSuite) TestDeleteTrain() {
	if suite.skip {
		return
	}
	assert.NotPanics(suite.T(), func() { suite.Store.SaveAll(annons) }, "Save all should not panic.")
	assert.Equal(suite.T(), 3, countAnnons(suite.Session, suite.d))
	assert.NotPanics(suite.T(), func() { suite.Store.DeleteForTrain(trainID) }, "Delete should not panic.")
	assert.Equal(suite.T(), 1, countAnnons(suite.Session, suite.d))
}

func TestAnnonStoreTestSuite(t *testing.T) {
	suite.Run(t, new(AnnonStoreTestSuite))
}

func (suite *AnnonStoreTestSuite) TestRealtimeChanges() {
	if suite.skip {
		return
	}
	a1 := Announcement{ID: "a1", Time: []time.Time{time.Now().Add(time.Minute * 3), time.Now().Add(time.Minute * 10)}}
	a2 := Announcement{ID: "a2", Time: []time.Time{time.Now().Add(time.Minute * -10), time.Now().Add(time.Minute * 4)}}
	a3 := Announcement{ID: "a3", Time: []time.Time{time.Now().Add(time.Minute * -10), time.Now().Add(time.Minute * -2)}}
	a4 := Announcement{ID: "a4", Time: []time.Time{time.Now().Add(time.Minute * 13), time.Now().Add(time.Minute * 10)}}
	a5 := Announcement{ID: "a5", Time: []time.Time{time.Now().Add(time.Minute * 3), time.Now().Add(time.Minute * 8)}}
	a6 := Announcement{ID: "a6", Time: []time.Time{time.Now().Add(time.Minute * 9), time.Now().Add(time.Minute * 10)}}
	as := []*Announcement{&a1, &a2, &a3, &a4}
	suite.Store.SaveAll(&as)
	checkpoints := 0
	go func() {
		next, old, err := suite.Store.NextChange()
		assert.NoError(suite.T(), err)
		assert.Nil(suite.T(), old)
		assert.Equal(suite.T(), "a5", next.ID)
		checkpoints++
		next, old, err = suite.Store.NextChange()
		assert.NoError(suite.T(), err)
		assert.Nil(suite.T(), old)
		assert.Equal(suite.T(), "a6", next.ID)
		checkpoints++
	}()
	time.Sleep(time.Millisecond * 300)
	suite.Store.Save(&a5)
	suite.Store.Save(&a6)
	assert.Equal(suite.T(), 2, checkpoints)

}

var annons = &[]*Announcement{
	&Announcement{
		Time:     []time.Time{time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)},
		First:    time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
		Last:     time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone),
		Type:     "arrival",
		Priority: P1,
		TrainID:  trainID,
		Text: &Text{
			HumanText: "I'm a human",
			HTMLText:  "<tag>I'm HTML</tag>",
			TtsText:   "I'm a TTS",
		},
	},
	&Announcement{
		Time:     []time.Time{time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone)},
		First:    time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone),
		Last:     time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone),
		Type:     "arrival",
		TrainID:  trainID,
		Priority: P1,
		Text: &Text{
			HumanText: "I'm a second human",
			HTMLText:  "<tag>I'm a second HTML</tag>",
			TtsText:   "I'm a second TTS",
		},
	},
	&Announcement{
		Time:    []time.Time{time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone)},
		First:   time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone),
		Last:    time.Date(2016, time.April, 10, 15, 25, 0, 0, config.Timezone),
		Type:    "arrival",
		TrainID: "another",
		Text: &Text{
			HumanText: "I'm a second human",
			HTMLText:  "<tag>I'm a second HTML</tag>",
			TtsText:   "I'm a second TTS",
		},
	},
}

func cleanupAnnons(session *r.Session, d db.Database) {
	r.DB(d).Table(db.AnnonTable).Delete().RunWrite(session)
}

func countAnnons(session *r.Session, d db.Database) int {
	cur, err := r.DB(d).Table(db.AnnonTable).Count().Run(session)
	if err != nil {
		panic(err)
	}
	defer cur.Close()
	var c int
	err = cur.One(&c)
	if err != nil {
		panic(err)
	}
	return c
}
