package main

import (
	"bufio"
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

	var format alsa.Format
	switch os.Args[2] {
	case "1":
		format = alsa.FormatU16LE
	case "2":
		format = alsa.FormatU16BE
	case "3":
		format = alsa.FormatS16LE
	case "4":
		format = alsa.FormatS16BE
	}

	var device *alsa.PlaybackDevice
	if device, err = alsa.NewPlaybackDevice("PCM", 1, format, 22050, alsa.BufferParams{BufferFrames: 10, PeriodFrames: 4, Periods: 2}); err != nil {
		fmt.Printf("Could not create device. %v", err)
		os.Exit(1)
	}

	var read int
	var wrote int
	defer device.Close()
	buf := make([]byte, 1024)
	for {
		if read, err = r.Read(buf); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Write error : %v\n", err)
			os.Exit(3)
		}
		if read == 0 {
			break
		}

		if wrote, err = device.Write(buf); err != nil {
			fmt.Fprintf(os.Stderr, "Write error : %v\n", err)
			os.Exit(4)
		}
		if wrote != read {
			fmt.Fprintf(os.Stderr, "Did not write whole buffer (Wrote %v, expected %v)\n", wrote, read)
		}
	}

}
