package annon

import (
	"testing"

	"golang.org/x/text/language"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func (suite *UtilsTestSuite) SetupSuite() {
}

func (suite *UtilsTestSuite) TestLocaleFormat() {
	assert.Equal(suite.T(), "15:04", TimeFormat(language.Polish))
	assert.Equal(suite.T(), "03:04pm", TimeFormat(language.English))
	assert.Equal(suite.T(), "15:04", TimeFormat(language.Russian))
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
