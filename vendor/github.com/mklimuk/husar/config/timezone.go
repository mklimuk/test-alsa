package config

import "time"

//Timezone is a default timezone for Husar
var Timezone, _ = time.LoadLocation("Europe/Warsaw")
