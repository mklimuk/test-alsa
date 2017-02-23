package config

import "time"

//Conf holds application's configuration
type Conf struct {
	GPIO  GPIOConf       `yaml:"gpio"`
	Reg   []RegistryConf `yaml:"registries"`
	Audio AudioConf      `yaml:"audio" json:"audio"`

	APIHostname string `yaml:"apiHost"`
	APIPort     string `yaml:"apiPort"`
	LogFile     string `yaml:"log"` //	/var/log/neo9/

	StateUpdateFreq  time.Duration `yaml:"stateUpdateFreq"`  //	5,
	RegistryBindFreq time.Duration `yaml:"registryBindFreq"` //	5,
}

//RegistryConf holds configuration o a single registry endpoint
type RegistryConf struct {
	Host string `yaml:"hostname" json:"hostname"` //	0.0.0.0
	Port string `yaml:"port"`                     //	7081
	IP   string `json:"ip"`
}

//AudioConf holds audio-related configuration
type AudioConf struct {
	Mixer         string   `yaml:"mixer"`              //	Master
	Zones         []string `yaml:"zones" json:"zones"` // default playback zones
	SetVolumePath string   `yaml:"setVol"`
	GetVolumePath string   `yaml:"getVol"`
	DeviceBuffer  int      `yaml:"deviceBuffer"`
	PeriodFrames  int      `yaml:"periodFrames"`
	Periods       int      `yaml:"periods"`
	ReadBuffer    int      `yaml:"readBuffer"` //in websocket frames
}

//GPIOConf holds I/O pin mappings and related info
type GPIOConf struct {
	Amp   AmpConf   `yaml:"amp"`
	Alarm AlarmConf `yaml:"alarm"`
}

//AmpConf holds I/O settings related to the amplifier
type AmpConf struct {
	OnOffTimeout        time.Duration `yaml:"onOffTimeout"`        //	10,
	OnOffPin            int           `yaml:"onOffPin"`            //	29
	PowerStatusPin      int           `yaml:"powerStatusPin"`      // 25
	ChannelAvailablePin int           `yaml:"channelAvailablePin"` // 27
}

//AlarmConf holds I/O settings related to the alarm logic
type AlarmConf struct {
	AlarmPin int `yaml:"alarmPin"` // 28
}
