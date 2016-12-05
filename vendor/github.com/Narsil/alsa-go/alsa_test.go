package alsa

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	handle := New()
	err := handle.Open("default", StreamTypePlayback, ModeBlock)
	if err != nil {
		t.Fatalf("Open failed. %s", err)
	}

	handle.SampleFormat = SampleFormatS16LE
	handle.SampleRate = 44100
	handle.Channels = 2
	err = handle.ApplyHwParams()
	if err != nil {
		t.Fatalf("SetHwParams failed. %s", err)
	}

	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	wrote, err := handle.Write(buf)
	if err != nil {
		t.Fatalf("Writei failed. %s", err)
	}
	if wrote != 20 {
		t.Errorf("Did not write all the buffer")
	}

	handle.Close()
}

func TestCapture(t *testing.T) {
	handle := New()
	err := handle.Open("default", StreamTypeCapture, ModeBlock)
	if err != nil {
		t.Fatalf("Open failed. %s", err)
	}

	handle.SampleFormat = SampleFormatU8
	handle.SampleRate = 8000
	handle.Channels = 1
	err = handle.ApplyHwParams()
	if err != nil {
		t.Fatalf("SetHwParams failed. %s", err)
	}

	buflen := int(1024)
	buf := make([]byte, buflen)
	n, err := handle.Read(buf)
	if err != nil {
		t.Fatalf("Writei failed. %s", err)
	}
	if len(buf) != n {
		t.Errorf("Could not read all data, Read %s (expected %s)", n, buflen)
	}
	handle.Close()

}

func BenchmarkRead(b *testing.B) {
	handle := New()
	err := handle.Open("default", StreamTypeCapture, ModeBlock)
	if err != nil {
		b.FailNow()
	}

	handle.SampleFormat = SampleFormatU8
	handle.SampleRate = 8000
	handle.Channels = 1
	err = handle.ApplyHwParams()
	if err != nil {
		b.FailNow()
	}

	buf := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_, err := handle.Read(buf)
		if err != nil {
			b.FailNow()
		}
	}
	handle.Close()
}

func BenchmarkWrite(b *testing.B) {
	handle := New()
	err := handle.Open("default", StreamTypePlayback, ModeBlock)
	if err != nil {
		b.FailNow()
	}

	handle.SampleFormat = SampleFormatU8
	handle.SampleRate = 8000
	handle.Channels = 1
	err = handle.ApplyHwParams()
	if err != nil {
		b.FailNow()
	}
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for i := 0; i < b.N; i++ {
		handle.Write(buf)
	}
	handle.Close()
}

func ExampleHandle_Read() {
	handle := New()
	err := handle.Open("default", StreamTypeCapture, ModeBlock)
	if err != nil {
		fmt.Printf("Open failed. %s", err)
	}

	handle.SampleFormat = SampleFormatU8
	handle.SampleRate = 8000
	handle.Channels = 1
	err = handle.ApplyHwParams()
	if err != nil {
		fmt.Printf("SetHwParams failed. %s", err)
	}

	buf := make([]byte, 1024)
	n, err := handle.Read(buf)
	if err != nil {
		fmt.Printf("Read failed. %s", err)
	}
	if n != len(buf) {
		fmt.Printf("Could not read all data (Read %i, expected %i)", n, len(buf))
	}
	handle.Close()
}

func ExampleHandle_Write() {
	handle := New()
	err := handle.Open("default", StreamTypePlayback, ModeBlock)
	if err != nil {
		fmt.Printf("Open failed. %s", err)
	}
	handle.SampleFormat = SampleFormatU8
	handle.SampleRate = 8000
	handle.Channels = 1
	err = handle.ApplyHwParams()
	if err != nil {
		fmt.Printf("SetHwParams failed. %s", err)
	}
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	n, err := handle.Write(buf)
	if err != nil {
		fmt.Printf("Write failed %s", err)
	}
	if n != len(buf) {
		fmt.Printf("Did not write all data (Wrote %s, expected %s)", n, len(buf))
	}
	handle.Close()
}
