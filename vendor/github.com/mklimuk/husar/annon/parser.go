package annon

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

/*Parse parses templates definitions from a configuration file*/
func Parse(path []string) CatalogTemplates {
	var cat []CatalogTemplate
	for _, p := range path {
		cat = append(cat, parseFile(p)...)
	}
	return CatalogTemplates(cat)
}

func parseFile(path string) []CatalogTemplate {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"file":  path,
		}).Panicln("Could not read the templates file")
	}
	var cat CatalogTemplates
	if err = yaml.Unmarshal(file, &cat); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"file":  path,
		}).Panicln("Could not parse the configuration file")
	}
	return cat
}
