package main

/*
MIT License

Copyright (c) 2018 Emmanuel Arias Soto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/eariassoto/mjpegstreamer"
)

func startSourceStream(dataChan chan<- []byte) {
	for {
		command := exec.Command("avconv",
			"-f", "x11grab",
			"-s", "1920:1080",
			"-i", ":0.0",
			"-s", "1280:720",
			"-f", "image2",
			"-frames", "1",
			"-",
		)
		errorBuffer := new(bytes.Buffer)
		command.Stderr = errorBuffer
		data, cmdErr := command.Output()

		if cmdErr != nil {
			log.Printf(cmdErr.Error())
			log.Printf("Error buffer: %s", errorBuffer)
			continue
		}
		dataChan <- data
	}
}

func main() {
	//host := "127.0.0.1"
	//port := 4000

	streamer := mjpegstreamer.NewMjpegStreamer()
	sourceChan := make(chan []byte)

	streamer.AddStream("desktop", sourceChan)
	go startSourceStream(sourceChan)
	streamer.StartStream()
}
