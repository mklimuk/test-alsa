package l10n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/text/language"
)

type MapperTestSuite struct {
	suite.Suite
}

func (suite *MapperTestSuite) TestRomanMapper() {
	assert.Equal(suite.T(), 1, RomanToInt("I"))
	assert.Equal(suite.T(), 4, RomanToInt("IV"))
	assert.Equal(suite.T(), 8, RomanToInt("VIII"))
	assert.Equal(suite.T(), 10, RomanToInt("X"))
}

func (suite *MapperTestSuite) TestNumberToTextMapper() {
	assert.Equal(suite.T(), "pierwszy", NumToText(1, Genitive, language.Polish))
	assert.Equal(suite.T(), "trzeciego", NumToText(3, Adjective, language.Polish))
	assert.Equal(suite.T(), "dziesiątym", NumToText(10, Locative, language.Polish))
}

func (suite *MapperTestSuite) TestMetaDictionary() {
	assert.Equal(suite.T(), "Kolei Mazowieckich", FromMetaDictionary("KM", Locative, language.Polish))
	assert.Equal(suite.T(), "Osobowy", FromMetaDictionary("Os", Genitive, language.Polish))
	assert.Equal(suite.T(), "", FromMetaDictionary("IC", Locative, language.Polish))
}

func TestMapperTestSuite(t *testing.T) {
	suite.Run(t, new(MapperTestSuite))
}
