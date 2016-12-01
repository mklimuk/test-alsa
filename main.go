package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"

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
	bFrames, _ := strconv.Atoi(os.Args[3])
	pFrames, _ := strconv.Atoi(os.Args[4])
	periods, _ := strconv.Atoi(os.Args[5])
	if device, err = alsa.NewPlaybackDevice(os.Args[2], 1, alsa.FormatS16LE, 22050, alsa.BufferParams{BufferFrames: bFrames, PeriodFrames: pFrames, Periods: periods}); err != nil {
		fmt.Printf("Could not create device. %v", err)
		os.Exit(1)
	}

	var read int
	var wrote int
	defer device.Close()

	bufSize, _ := strconv.Atoi(os.Args[6])
	buf := make([]byte, bufSize)
	buf16 := make([]int16, bufSize/2)

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
