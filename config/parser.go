package config

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

/*
Parse parses the configuration file into Config
*/
func Parse(path string) *Conf {
	var file []byte
	var err error
	if file, err = ioutil.ReadFile(path); err != nil {
		log.WithFields(log.Fields{
			"file": path,
		}).Panicln("Could not read the configuration file")
	}
	var conf Conf
	if err = yaml.Unmarshal(file, &conf); err != nil {
		log.WithFields(log.Fields{
			"file": path,
		}).Panicln("Could not parse the configuration file")
	}
	return &conf
}
