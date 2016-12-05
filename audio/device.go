package audio

//BufferParams is a copy of Alsa configuration parameters present in audio package for isolation purposes
//(to allow testability of audio package without cgo)
type BufferParams struct {
	BufferFrames int
	PeriodFrames int
	Periods      int
}

//PlaybackDevice is an interface wrapper over goalsa PlaybackDevice.
//It is defined for testing convienience (mocking and cgo independence)
type PlaybackDevice interface {
	Write(buffer interface{}) (samples int, err error)
	Close()
}
