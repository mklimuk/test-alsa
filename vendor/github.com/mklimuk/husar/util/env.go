package util

import "os"

const versionFile = "/var/husar/version.yml"

//GetEnv returns environment variable if it is set or defaultValue otherwise
func GetEnv(key string, defaultValue string) (val string) {
	if val = os.Getenv(key); val == "" {
		return defaultValue
	}
	return val
}
