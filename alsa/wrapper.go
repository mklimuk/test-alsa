package alsa

import (
	goalsa "github.com/mklimuk/goalsa"
	"github.com/mklimuk/test-alsa/audio"
)

const deviceName = "sysdefault"
const format = goalsa.FormatS16LE

//Factory implements audio.DeviceFactory interface
type Factory struct {
}

//New wraps goalsa PlaybackDevice constructor for testing convienience
func (f *Factory) New(sampleRate int, channels int, bp *audio.BufferParams) (audio.PlaybackDevice, error) {
	return goalsa.NewPlaybackDevice(deviceName, channels, format, sampleRate, goalsa.BufferParams{BufferFrames: bp.BufferFrames, PeriodFrames: bp.PeriodFrames, Periods: bp.Periods})
}
