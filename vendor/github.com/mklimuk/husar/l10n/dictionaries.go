package l10n

import "golang.org/x/text/language"

// Category represents language constructs category to be used in l10n dictionaries
type Category int

// Constants representing different language element categories
const (
	Genitive Category = iota
	Adjective
	Locative
	Cardinal
	Ordinal
)

var roman map[string]int
var numbers map[language.Tag]map[int]map[Category]string
var dictionary map[language.Tag]map[string]map[Category]string

func init() {
	roman = map[string]int{
		"I":    1,
		"II":   2,
		"III":  3,
		"IV":   4,
		"V":    5,
		"VI":   6,
		"VII":  7,
		"VIII": 8,
		"IX":   9,
		"X":    10,
	}
	numbers = map[language.Tag]map[int]map[Category]string{
		language.Polish: map[int]map[Category]string{
			1: map[Category]string{
				Genitive:  "pierwszy",
				Adjective: "pierwszego",
				Locative:  "pierwszym",
			},
			2: map[Category]string{
				Genitive:  "drugi",
				Adjective: "drugiego",
				Locative:  "drugim",
			},
			3: map[Category]string{
				Genitive:  "trzeci",
				Adjective: "trzeciego",
				Locative:  "trzecim",
			},
			4: map[Category]string{
				Genitive:  "czwarty",
				Adjective: "czwartego",
				Locative:  "czwartym",
			},
			5: map[Category]string{
				Genitive:  "piąty",
				Adjective: "piątego",
				Locative:  "piątym",
			},
			6: map[Category]string{
				Genitive:  "szósty",
				Adjective: "szóstego",
				Locative:  "szóstym",
			},
			7: map[Category]string{
				Genitive:  "siódmy",
				Adjective: "siódmego",
				Locative:  "siódmym",
			},
			8: map[Category]string{
				Genitive:  "ósmy",
				Adjective: "ósmego",
				Locative:  "ósmym",
			},
			9: map[Category]string{
				Genitive:  "dziewiąty",
				Adjective: "dziewiątego",
				Locative:  "dziewiątym",
			},
			10: map[Category]string{
				Genitive:  "dziesiąty",
				Adjective: "dziesiątego",
				Locative:  "dziesiątym",
			},
			11: map[Category]string{
				Genitive:  "jedenasty",
				Adjective: "jedenastego",
				Locative:  "jedenastym",
			},
			12: map[Category]string{
				Genitive:  "dwunasty",
				Adjective: "dwunastego",
				Locative:  "dwunastym",
			},
			13: map[Category]string{
				Genitive:  "trzynasty",
				Adjective: "trzynastego",
				Locative:  "trzynastym",
			},
			14: map[Category]string{
				Genitive:  "czternasty",
				Adjective: "czternastego",
				Locative:  "czternastym",
			},
			15: map[Category]string{
				Genitive:  "piętnasty",
				Adjective: "piętnastego",
				Locative:  "piętnastym",
			},
		},
		language.English: map[int]map[Category]string{
			1: map[Category]string{
				Genitive:  "one",
				Adjective: "one",
				Locative:  "one",
			},
			2: map[Category]string{
				Genitive:  "two",
				Adjective: "two",
				Locative:  "two",
			},
			3: map[Category]string{
				Genitive:  "three",
				Adjective: "three",
				Locative:  "three",
			},
			4: map[Category]string{
				Genitive:  "four",
				Adjective: "four",
				Locative:  "four",
			},
			5: map[Category]string{
				Genitive:  "five",
				Adjective: "five",
				Locative:  "five",
			},
			6: map[Category]string{
				Genitive:  "six",
				Adjective: "six",
				Locative:  "six",
			},
			7: map[Category]string{
				Genitive:  "seven",
				Adjective: "seven",
				Locative:  "seven",
			},
			8: map[Category]string{
				Genitive:  "eight",
				Adjective: "eight",
				Locative:  "eight",
			},
			9: map[Category]string{
				Genitive:  "nine",
				Adjective: "nine",
				Locative:  "nine",
			},
			10: map[Category]string{
				Genitive:  "ten",
				Adjective: "ten",
				Locative:  "ten",
			},
			11: map[Category]string{
				Genitive:  "eleven",
				Adjective: "eleven",
				Locative:  "eleven",
			},
			12: map[Category]string{
				Genitive:  "twelve",
				Adjective: "twelve",
				Locative:  "twelve",
			},
			13: map[Category]string{
				Genitive:  "thirteen",
				Adjective: "thirteen",
				Locative:  "thirteen",
			},
			14: map[Category]string{
				Genitive:  "fourteen",
				Adjective: "fourteen",
				Locative:  "fourteen",
			},
			15: map[Category]string{
				Genitive:  "fifteen",
				Adjective: "fifteen",
				Locative:  "fifteen",
			},
		},
	}
	dictionary = map[language.Tag]map[string]map[Category]string{
		language.Polish: map[string]map[Category]string{
			"KM": map[Category]string{
				Genitive: "Koleje Mazowieckie",
				Locative: "Kolei Mazowieckich",
			},
			"TLK": map[Category]string{
				Genitive: "TLK",
				Locative: "TLK",
			},
			"Os": map[Category]string{
				Genitive: "Osobowy",
				Locative: "Osobowym",
			},
			"SDM": map[Category]string{
				Genitive: "Specjalny",
				Locative: "Specjalnym",
			},
			"S": map[Category]string{
				Genitive: "Specjalny",
				Locative: "Specjalnym",
			},
		},
		language.English: map[string]map[Category]string{
			"KM": map[Category]string{
				Genitive: "Koleje Mazowieckie",
				Locative: "Koleje Mazowieckie",
			},
			"TLK": map[Category]string{
				Genitive: "TLK",
				Locative: "TLK",
			},
			"Os": map[Category]string{
				Genitive: "",
				Locative: "",
			},
			"SDM": map[Category]string{
				Genitive: "",
				Locative: "",
			},
			"S": map[Category]string{
				Genitive: "",
				Locative: "",
			},
		},
	}
}
