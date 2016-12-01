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
	if device, err = alsa.NewPlaybackDevice(os.Args[3], 1, alsa.FormatU16LE, 22050, alsa.BufferParams{BufferFrames: 10, PeriodFrames: 4, Periods: 2}); err != nil {
		fmt.Printf("Could not create device. %v", err)
		os.Exit(1)
	}

	var read int
	var wrote int
	defer device.Close()
	buf := make([]byte, 1024)
	buf16 := make([]int16, 512)

	for {
		if read, err = r.Read(buf); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Write error : %v\n", err)
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
		if wrote != read {
			fmt.Fprintf(os.Stderr, "Did not write whole buffer (Wrote %v, expected %v)\n", wrote, read)
		}
	}

}

func convertBuffers(buf []byte, buf16 []int16) {
	for i := range buf {
		fmt.Println(i)
		// assuming little endian
		buf16[i] = int16(binary.LittleEndian.Uint16(buf[i*sizeofInt16 : (i+1)*sizeofInt16]))
	}
}
