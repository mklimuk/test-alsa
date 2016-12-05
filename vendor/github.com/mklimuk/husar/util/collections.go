package util

import "reflect"

/*SliceContains checks if a given element can be found in a slice of strings.
It is designed work with any type that can be converted to string via reflection. */
func SliceContains(list []string, toFind interface{}) bool {
	s := reflect.ValueOf(toFind)
	for _, element := range list {
		if element == s.String() {
			return true
		}
	}
	return false
}
