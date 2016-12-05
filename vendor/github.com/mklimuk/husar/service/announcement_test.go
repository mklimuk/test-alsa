package service

import (
	"errors"
	"testing"
	"time"

	"github.com/mklimuk/husar/annon"
	"github.com/mklimuk/husar/audio"
	"github.com/mklimuk/husar/config"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (suite *ServiceTestSuite) TestGenError() {
	t := &train.Train{ID: "test"}
	trainStore := &test.TrainStoreMock{}
	annonStore := &annon.AnnonStoreMock{}
	aud := &test.AudioCatalogMock{}
	gen := &annon.GeneratorMock{}
	cat := &annon.CatalogMock{}
	current := new([]*annon.Announcement)
	annonStore.On("GetForTrain", t.ID).Return(current, nil)
	gen.On("ForTrain", t, t, current, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*annon.Catalog")).Return([]*annon.Announcement{}, []*annon.Announcement{}, errors.New("Test error"))
	aud.On("Generate", mock.Anything).Return("ID", true, 10, false, nil)
	g := annon.Generator(gen)
	as := annon.Store(annonStore)
	ts := train.Store(trainStore)
	c := annon.Catalog(cat)
	a := audio.Catalog(aud)
	sync := NewAnnouncement(&g, &as, &ts, &c, &a)
	err := sync.(*announcement).ProcessTrain(t, t)
	assert.NotNil(suite.T(), err)
	trainStore.AssertExpectations(suite.T())
	annonStore.AssertExpectations(suite.T())
	gen.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCallProcessing() {
	t := &train.Train{ID: "test"}
	text := annon.Text{TtsText: "test"}
	a := &[]*annon.Announcement{&annon.Announcement{Text: &text}}
	trainStore := &test.TrainStoreMock{}
	annonStore := &annon.AnnonStoreMock{}
	gen := &annon.GeneratorMock{}
	aud := &test.AudioCatalogMock{}
	cat := &annon.CatalogMock{}
	current := new([]*annon.Announcement)
	args := annonStore.On("SaveAll", a).Return(nil).Arguments
	annonStore.On("GetForTrain", t.ID).Return(current, nil)
	gen.On("ForTrain", t, t, current, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*annon.Catalog")).Return(*a, []*annon.Announcement{}, nil)
	aud.On("Generate", mock.Anything).Return("ID", true, 10, false, nil)
	g := annon.Generator(gen)
	as := annon.Store(annonStore)
	ts := train.Store(trainStore)
	c := annon.Catalog(cat)
	au := audio.Catalog(aud)
	sync := NewAnnouncement(&g, &as, &ts, &c, &au)
	err := sync.(*announcement).ProcessTrain(t, t)
	assert.Nil(suite.T(), err)
	annons := args.Get(0).(*[]*annon.Announcement)
	assert.NotNil(suite.T(), (*annons)[0].Audio)
	trainStore.AssertExpectations(suite.T())
	annonStore.AssertExpectations(suite.T())
	gen.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestDelete() {
	s := announcement{}
	assert.NotPanics(suite.T(), func() {
		s.ProcessRealtimeChange(nil, nil)
	})
	assert.NotPanics(suite.T(), func() {
		s.ProcessSettingsChange(nil, nil)
	})
}

func (suite *ServiceTestSuite) TestGenerateFromCatalog() {
	cat := &annon.CatalogMock{}
	c := annon.Catalog(cat)
	cat.On("Get", "abcd").Return(nil, errors.New("Test error"))
	s := NewAnnouncement(nil, nil, nil, &c, nil)
	err := s.GenerateFromCatalog("abcd", "zone1", nil, nil, 10, 10, nil)
	assert.Error(suite.T(), err)
	start := time.Date(2016, time.April, 10, 15, 21, 0, 0, config.Timezone)
	err = s.GenerateFromCatalog("abcd", "zone1", &start, nil, 10, 10, nil)
	assert.Error(suite.T(), err)
}

func (suite *ServiceTestSuite) TestGetTrainsInRange() {
	in1 := time.Date(2016, time.April, 10, 15, 2, 0, 0, config.Timezone)
	in2 := time.Date(2016, time.April, 10, 14, 53, 0, 0, config.Timezone)
	in2delay := time.Date(2016, time.April, 10, 15, 8, 0, 0, config.Timezone) //15 minutes
	in3 := time.Date(2016, time.April, 10, 15, 52, 0, 0, config.Timezone)
	in3delay := time.Date(2016, time.April, 10, 15, 57, 0, 0, config.Timezone) //5
	in4 := time.Date(2016, time.April, 10, 16, 55, 0, 0, config.Timezone)
	in4delay := time.Date(2016, time.April, 10, 17, 5, 0, 0, config.Timezone) // 5 minutes
	in5 := time.Date(2016, time.April, 10, 16, 12, 0, 0, config.Timezone)
	in6 := time.Date(2016, time.April, 10, 16, 34, 0, 0, config.Timezone)
	in6delay := time.Date(2016, time.April, 10, 16, 44, 0, 0, config.Timezone) //10 minutes

	trains := []train.Train{
		train.Train{ID: "t2", FirstEvent: &in5, FirstLiveEvent: &in5},
		train.Train{ID: "t1", FirstEvent: &in6, Delay: 10, FirstLiveEvent: &in6delay},
		train.Train{ID: "t4", FirstEvent: &in2, Delay: 15, FirstLiveEvent: &in2delay},
		train.Train{ID: "t5", FirstEvent: &in3, Delay: 5, FirstLiveEvent: &in3delay},
		train.Train{ID: "t6", FirstEvent: &in4, Delay: 5, FirstLiveEvent: &in4delay},
		train.Train{ID: "t3", FirstEvent: &in1, FirstLiveEvent: &in1},
	}
	settings := []train.Settings{
		train.Settings{ID: "t2"},
		train.Settings{ID: "t3"},
		train.Settings{ID: "t5"},
	}

	start := time.Date(2016, time.April, 10, 15, 0, 0, 0, config.Timezone)
	end := time.Date(2016, time.April, 10, 17, 0, 0, 0, config.Timezone)

	store := &test.TrainStoreMock{}
	store.On("GetBetweenWithRealtime", &start, &end).Return(trains[1:], []string{"t1", "t4", "t5", "t6", "t3"}, nil)
	store.AssertNotCalled(suite.T(), "DeleteForTrain", mock.AnythingOfType("string"))
	store.On("GetAllSettings", []string{"t1", "t4", "t5", "t6", "t3"}).Return(&settings, nil)
	anno := &annon.AnnonStoreMock{}
	aud := &test.AudioCatalogMock{}
	tex := annon.Text{TtsText: "test"}
	a := []*annon.Announcement{&annon.Announcement{Text: &tex}}
	anno.On("GetForTrain", mock.Anything).Times(5).Return(nil, nil)
	anno.On("SaveAll", &a).Times(5).Return(nil)
	aud.On("Generate", mock.Anything).Return("ID", true, 10, false, nil)
	gen := &annon.GeneratorMock{}
	gen.On("ForTrain", mock.AnythingOfType("*train.Train"), (*train.Train)(nil), (*[]*annon.Announcement)(nil), mock.AnythingOfType("*time.Time"), (*annon.Catalog)(nil)).Times(5).Return(a, []*annon.Announcement{}, nil)
	g := annon.Generator(gen)
	st := train.Store(store)
	an := annon.Store(anno)
	au := audio.Catalog(aud)
	s := NewAnnouncement(&g, &an, &st, nil, &au)
	s.GenerateForPeriod(&start, &end)
	store.AssertExpectations(suite.T())
	anno.AssertExpectations(suite.T())
	gen.AssertExpectations(suite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
