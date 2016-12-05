package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/*
Parse parses the configuration file into an arbitrary configuration
*/
func Parse(path *string, config interface{}) (err error) {
	var file []byte
	if file, err = ioutil.ReadFile(*path); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, config); err != nil {
		return err
	}
	return
}
