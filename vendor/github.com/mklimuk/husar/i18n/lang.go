package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var supported = []language.Tag{
	language.Polish,
	language.English,
}

var matcher = language.NewMatcher(supported)

//DefaultLang defines default system locale
var DefaultLang = language.Polish

//GetLang returns language corresponding to a given string locale. Returns default language if not specified.
func GetLang(locale string) (l language.Tag) {
	var err error
	if l, err = language.Parse(locale); err != nil {
		return DefaultLang
	}
	l, _, _ = matcher.Match(l)
	return l
}

//GetSupported returns a map of all supported languages with their names in Polish.
func GetSupported() map[string]string {
	res := make(map[string]string)
	pl := display.Polish.Tags()
	for _, l := range supported {
		res[l.String()] = pl.Name(l)
	}
	return res
}
