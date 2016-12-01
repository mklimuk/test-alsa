package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	alsa "github.com/cocoonlife/goalsa"
)

func main() {
	file := os.Args[1]
	var f *os.File
	var err error
	if f, err = os.Open(file); err != nil {
		fmt.Printf("Could not open file %s", file)
		os.Exit(1)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	var device *alsa.PlaybackDevice
	if device, err = alsa.NewPlaybackDevice(os.Args[3], 1, alsa.FormatS16LE, 22050, alsa.BufferParams{BufferFrames: 1024, PeriodFrames: 256, Periods: 4}); err != nil {
		fmt.Printf("Could not create device. %v", err)
		os.Exit(1)
	}

	var read int
	var wrote int
	defer device.Close()
	buf := make([]byte, 2048)
	buf16 := make([]int16, 1024)

	for {
		if read, err = r.Read(buf); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Read error : %v\n", err)
			os.Exit(3)
		}
		if read == 0 {
			break
		}

		convertBuffers(buf, buf16)
		if wrote, err = device.Write(buf16); err != nil {
			fmt.Fprintf(os.Stderr, "Write error : %v\n", err)
			os.Exit(4)
		}

		if wrote != read/2 {
			fmt.Printf("Read %d bytes, wrote %d frames", read, wrote)
		}

	}

}

func convertBuffers(buf []byte, buf16 []int16) {
	for i := 0; i < len(buf16); i++ {
		// assuming little endian
		buf16[i] = int16(binary.LittleEndian.Uint16(buf[i*2 : (i+1)*2]))
	}
}
