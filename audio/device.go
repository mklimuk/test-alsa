package audio

import (
	"io"

	log "github.com/Sirupsen/logrus"
)

const sampleSizeBytes = 2

//BufferParams is a copy of Alsa configuration parameters present in audio package for isolation purposes
//(to allow testability of audio package without cgo)
type BufferParams struct {
	BufferFrames int
	PeriodFrames int
	Periods      int
}

//RawDevice is an interface wrapper over goalsa PlaybackDevice.
//It is defined for testing convienience (mocking and cgo independence)
type RawDevice interface {
	Write(buffer interface{}) (samples int, err error)
	Close()
}

//PlaybackDevice is responsible for sending data to the audio device
type PlaybackDevice interface {
	WriteSync(reader io.Reader) error
	WriteAsync(buffer chan []int16) chan error
	FramesWrote() int
	Close()
}

//DeviceFactory provides new initialized playback devices
type DeviceFactory interface {
	New(sampleRate int, channels int, bp *BufferParams) (PlaybackDevice, error)
}

type dev struct {
	bufferSize  int
	framesWrote int
	ctrl        chan bool
	errors      chan error
	raw         RawDevice
}

//NewPlaybackDevice is the audio device constructor
func NewPlaybackDevice(raw RawDevice, bufferSize int) PlaybackDevice {
	d := dev{bufferSize: bufferSize, raw: raw}
	return &d
}

func (d *dev) WriteSync(reader io.Reader) error {

	buf := make([]byte, d.bufferSize)
	buf16 := make([]int16, d.bufferSize/sampleSizeBytes)

	var err error
	var read int
	var wrote int

	for {
		if read, err = reader.Read(buf); err != nil && err != io.EOF {
			return err
		}
		if read == 0 {
			break
		}
		convertBuffers(buf, buf16)
		if wrote, err = d.raw.Write(buf16); err != nil {
			log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "playFile"}).
				WithError(err).Error("Could not write buffer content to device")
			return err
		}
		d.framesWrote += wrote
	}
	return nil
}

func (d *dev) WriteAsync(buffer chan []int16) chan error {
	d.ctrl = make(chan bool)
	d.errors = make(chan error)
	go d.sendToDevice(buffer)
	return d.errors
}

func (d *dev) sendToDevice(buffer chan []int16) {
	var err error
	var wrote int
	var frame []int16
	var ok bool

	for {
		select {
		case frame, ok = <-buffer:
			if !ok {
				if log.GetLevel() >= log.InfoLevel {
					log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "sendToDevice"}).
						Info("Buffer channel is closed; aborting write routine")
				}
				return
			}
			if wrote, err = d.raw.Write(frame); err != nil {
				d.errors <- err
			}
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "sendToDevice", "wroteFrames": wrote}).
					Debug("Wrote read buffer to device")
			}
			d.framesWrote += wrote
		case <-d.ctrl:
			if log.GetLevel() >= log.InfoLevel {
				log.WithFields(log.Fields{"logger": "audio-endpoint.audio", "method": "sendToDevice"}).
					Info("Got a signal through the control channel; aborting send loop")
			}
			return
		}
	}
}

func (d *dev) Close() {
	if d.ctrl != nil {
		d.ctrl <- true
		close(d.ctrl)
	}
	if d.errors != nil {
		close(d.errors)
	}
	d.raw.Close()
}

func (d *dev) FramesWrote() int {
	return d.framesWrote
}
