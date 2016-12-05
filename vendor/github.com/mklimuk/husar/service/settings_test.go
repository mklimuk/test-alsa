package service

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mklimuk/husar/event"
	"github.com/mklimuk/husar/test"
	"github.com/mklimuk/husar/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SettingsTestSuite struct {
	suite.Suite
	b *event.Bus
}

func (suite *SettingsTestSuite) SetupSuite() {
	suite.b = event.New()
	log.SetLevel(log.DebugLevel)
}

func (suite *SettingsTestSuite) TestFirstCall() {
	st := test.TrainStoreMock{}
	st.On("SaveSettings", mock.AnythingOfType("*train.Settings")).Return(&train.Settings{}, nil)
	st.On("Get", "48355_20160727_90872859").Return(&train.Train{}, nil)
	st.On("GetAllSettings", []string{"48355_20160727_90872859"}).Return(&[]train.Settings{}, nil)
	t := timetable{&event.Bus{}, &st}
	t.updateSettings(`{"trainId":"48355_20160727_90872859","settings":null,"lang":{"en":{"enabled":true}}}`)
	settings := st.Calls[2].Arguments.Get(0).(*train.Settings)
	assert.Len(suite.T(), *settings.Lang, 1)
	assert.True(suite.T(), (*settings.Lang)["en"].Enabled)
}

func (suite *SettingsTestSuite) TestDisable() {
	st := test.TrainStoreMock{}
	st.On("SaveSettings", mock.AnythingOfType("*train.Settings")).Return(&train.Settings{}, nil)
	st.On("Get", "48355_20160727_90872859").Return(&train.Train{
		Settings: &train.Settings{
			Lang: &map[string]train.LanguageSettings{
				"en": train.LanguageSettings{
					Enabled: true,
				},
			},
		},
	}, nil)
	st.On("GetAllSettings", []string{"48355_20160727_90872859"}).Return(&[]train.Settings{}, nil)
	t := timetable{&event.Bus{}, &st}
	t.updateSettings(`{"trainId":"48355_20160727_90872859","settings":null,"lang":{"en":{"enabled":false}}}`)
	settings := st.Calls[2].Arguments.Get(0).(*train.Settings)
	assert.Len(suite.T(), *settings.Lang, 1)
	assert.False(suite.T(), (*settings.Lang)["en"].Enabled)
}

func TestSettingsTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsTestSuite))
}
