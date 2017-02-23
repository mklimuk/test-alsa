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
	var err error
	var dev *goalsa.PlaybackDevice
	if dev, err = goalsa.NewPlaybackDevice(deviceName, channels, format, sampleRate, goalsa.BufferParams{BufferFrames: bp.BufferFrames, PeriodFrames: bp.PeriodFrames, Periods: bp.Periods}); err != nil {
		return nil, err
	}
	return audio.NewPlaybackDevice(dev, bp.BufferFrames), nil
}
