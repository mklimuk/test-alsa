package annon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) TestParse() {
	path := "./templates.yml"
	templates := Parse([]string{path})
	assert.Len(suite.T(), templates, 10, "The file contains 10 templates")
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
