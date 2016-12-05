package device

import (
	"testing"

	"github.com/mklimuk/husar/event"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegistryTestSuite struct {
	suite.Suite
}

func (suite *RegistryTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)

}

func (suite *RegistryTestSuite) TestAddZone() {
	b := event.New()
	r := NewRegistry(b, 10)
	r.CreateZone("48355")
	assert.Len(suite.T(), r.(*reg).zones, 1)
	z := r.GetZone("48355")
	assert.NotNil(suite.T(), z)
	d := (*z).GetDevices()
	assert.NotNil(suite.T(), d)
}

func (suite *RegistryTestSuite) TestBindDevice() {
	b := event.New()
	r := NewRegistry(b, 10)
	d := Device{ID: "a", Bind: false, Zones: []string{"48355"}}
	r.CreateZone("48355")
	r.Bind(&d)
	assert.True(suite.T(), d.Bind)
	assert.Len(suite.T(), (*r.GetZone("48355")).GetDevices(), 1)
	r.AddToZone("48355", "a")
	err := r.Unregister("a")
	assert.NoError(suite.T(), err)
	dev, err := r.GetDevice("a")
	assert.Nil(suite.T(), dev)
	assert.Error(suite.T(), err)
}

func TestRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}
