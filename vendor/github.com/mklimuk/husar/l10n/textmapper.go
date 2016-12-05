package l10n

import "golang.org/x/text/language"

// NumToText returns a requested
func NumToText(num int, cat Category, lang language.Tag) string {
	return numbers[lang][num][cat]
}

// RomanToInt converts a roman number string into an integer (supports 1-10)
func RomanToInt(r string) int {
	return roman[r]
}

//FromMetaDictionary performs text translation from a general meta dictionary
func FromMetaDictionary(s string, cat Category, lang language.Tag) string {
	return dictionary[lang][s][cat]
}
