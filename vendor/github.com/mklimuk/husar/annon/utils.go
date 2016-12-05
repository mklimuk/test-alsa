package annon

import (
	"fmt"

	"golang.org/x/text/language"
)

//BuildID creates announcement ID based on it's type and trainId
func BuildID(trainID string, atype string, lang language.Tag) string {
	return fmt.Sprintf("%v_%v_%v", trainID, string(atype), lang)
}

/*MapFormat returns time format to be applied depending on locale*/
func TimeFormat(lang language.Tag) string {
	switch lang {
	case language.English:
		return "03:04pm"
	default:
		return "15:04"
	}
}

func CopyParams(params TemplateParams) TemplateParams {
	t := TemplateParams{}
	for i, p := range params {
		t[i] = p
	}
	return t
}
