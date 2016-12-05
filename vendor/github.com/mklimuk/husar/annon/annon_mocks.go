package annon

import (
	"time"

	"github.com/mklimuk/husar/train"
	"golang.org/x/text/language"

	"github.com/stretchr/testify/mock"
)

var prios = map[string]int{
	"arrival":   1,
	"departure": 1,
	"delay":     2,
	"setup":     2,
	"stop":      2,
	"info":      5,
}

//GeneratorMock is a mocked service.Generator interface
type GeneratorMock struct {
	mock.Mock
}

//ForTrain is a mocked method
func (m *GeneratorMock) ForTrain(t *train.Train, old *train.Train, currentAnnons *[]*Announcement, now *time.Time, cat *Catalog) (announcements []*Announcement, toDelete []*Announcement, err error) {
	args := m.Called(t, old, currentAnnons, now, cat)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Error(2)
		}
		return nil, args.Get(1).([]*Announcement), args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*Announcement), nil, args.Error(2)
	}
	return args.Get(0).([]*Announcement), args.Get(1).([]*Announcement), args.Error(2)
}

//Generate is a mocked method
func (m *GeneratorMock) Generate(t *train.Train, now *time.Time, temp *CatalogTemplate, p *Producer, lang language.Tag) (*Announcement, error) {
	args := m.Called(t, now, temp, p, lang)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Announcement), args.Error(1)
}

//Preview is a mocked method
func (m *GeneratorMock) Preview(tpl *CatalogTemplate, params TemplateParams) (*Text, error) {
	args := m.Called(tpl, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Text), args.Error(1)
}

//FromTemplate is a mocked method
func (m *GeneratorMock) FromTemplate(tpl *CatalogTemplate, startTime *time.Time, endTime *time.Time, dayInterval *time.Duration, nightInterval *time.Duration, now *time.Time, params TemplateParams, stationID int) (*Announcement, error) {
	args := m.Called(tpl, startTime, endTime, dayInterval, nightInterval, now, params, stationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Announcement), args.Error(1)
}

//ProducerMock is a mocked Producer interface
type ProducerMock struct {
	mock.Mock
}

//MapParams is a mocked method
func (m *ProducerMock) MapParams(params TemplateParams, lang language.Tag) TemplateParams {
	args := m.Called(params)
	return args.Get(0).(TemplateParams)
}

//Required is a mocked method
func (m *ProducerMock) Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error) {
	args := m.Called(t, old, reference)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

//GetTime is a mocked method
func (m *ProducerMock) GetTime(t *train.Train, now *time.Time) ([]time.Time, time.Time, time.Time) {
	args := m.Called(t, now)
	return args.Get(0).([]time.Time), args.Get(1).(time.Time), args.Get(2).(time.Time)
}

//BuildParams is a mocked method
func (m *ProducerMock) BuildParams(t *train.Train, lang language.Tag) (params *TemplateParams, ttsParams *TemplateParams, autoplay bool, err error) {
	args := m.Called(t)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil, args.Bool(2), args.Error(3)
		}
		return nil, args.Get(1).(*TemplateParams), args.Bool(2), args.Error(3)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*TemplateParams), nil, args.Bool(2), args.Error(3)
	}
	return args.Get(0).(*TemplateParams), args.Get(1).(*TemplateParams), args.Bool(2), args.Error(3)

}

//Name is a mocked method
func (m *ProducerMock) Name() string {
	args := m.Called()
	return args.String(0)
}

//ServiceMock is a mocked service.Service interface
type ServiceMock struct {
	mock.Mock
}

//ProcessRealtimeChange is a mocked method
func (s *ServiceMock) ProcessRealtimeChange(old, new *train.Realtime) {
	s.Called(old, new)
}

//ProcessSettingsChange is a mocked method
func (s *ServiceMock) ProcessSettingsChange(old, new *train.Settings) {
	s.Called(old, new)
}

//GenerateFromCatalog is a mocked method
func (s *ServiceMock) GenerateFromCatalog(ID string, title string, startTime *time.Time, endTime *time.Time, dayInterval int, nightInterval int, params *map[string]string) (err error) {
	args := s.Called(ID, title, startTime, endTime, dayInterval, nightInterval, params)
	return args.Error(0)
}

//GeneratePreview is a mocked method
func (s *ServiceMock) GeneratePreview(catalogID string, params *map[string]string) (txt *Text, duration int, file string, err error) {
	args := s.Called(catalogID, params)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.String(2), args.Error(3)
	}
	return args.Get(0).(*Text), args.Int(1), args.String(2), args.Error(3)
}

//GenerateForPeriod is a mocked method
func (s *ServiceMock) GenerateForPeriod(start, end *time.Time) {
	s.Called(start, end)
}

//GetCatalog is a mocked method
func (s *ServiceMock) GetCatalog() (tpl []CatalogTemplate, err error) {
	args := s.Called(tpl, err)
	// TODO wft?
	return nil, args.Error(1)
}

//CatalogMock is a mocked Catalog interface
type CatalogMock struct {
	mock.Mock
}

//Get is a mocked method
func (c *CatalogMock) Get(ID string) (*CatalogTemplate, error) {
	args := c.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CatalogTemplate), args.Error(1)
}

//Delete is a mocked method
func (c *CatalogMock) Delete(ID string) error {
	args := c.Called(ID)
	return args.Error(0)
}

//GetAll is a mocked method
func (c *CatalogMock) GetAll() (catalog map[string][]*CatalogTemplate, err error) {
	args := c.Called()
	return args.Get(0).(map[string][]*CatalogTemplate), args.Error(1)
}

//Save is a mocked method
func (c *CatalogMock) Save(template *CatalogTemplate) (*CatalogTemplate, error) {
	args := c.Called(template)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CatalogTemplate), args.Error(1)
}

//SaveAll is a mocked method
func (c *CatalogMock) SaveAll(tps *CatalogTemplates) (*CatalogTemplates, error) {
	args := c.Called(tps)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CatalogTemplates), args.Error(1)
}

//ByType is a mocked method
func (c *CatalogMock) ByType(t string) (*CatalogTemplate, error) {
	args := c.Called(t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CatalogTemplate), args.Error(1)
}
