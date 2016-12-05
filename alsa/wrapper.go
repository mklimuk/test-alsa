package alsa

import (
	goalsa "github.com/mklimuk/goalsa"
	"github.com/mklimuk/test-alsa/audio"
)

const channels = 1
const rate = 22050
const deviceName = "sysdefault"
const format = alsa.FormatS16LE

//NewPlaybackDevice wraps goalsa PlaybackDevice constructor for testing convienience
func NewPlaybackDevice(bp *audio.BufferParams) (p *audio.PlaybackDevice, err error) {
	return goalsa.NewPlaybackDevice(deviceName, channels, format, rate, goalsa.BufferParams{BufferFrames: bp.BufferFrames, PeriodFrames: bp.PeriodFrames, Periods: bp.Periods})
}
