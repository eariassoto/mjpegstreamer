package mjpegstreamer

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
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hybridgroup/mjpeg"
)

const (
	indexTemplate = `<!DOCTYPE html>
	<html>
	<head>
		<title>MJPEG data stream</title>
	</head>
	<body>
		<img src="/desktop">
	</body>
	</html>
	`
)

type MjpegStreamer struct {
	serveMux *http.ServeMux
	sources  map[string]*Stream
}

type Stream struct {
	Name        string
	sourceData  chan []byte
	done        chan bool
	mjpegstream *mjpeg.Stream
}

func (s *Stream) updateStream() {
	for {
		select {
		case data := <-s.sourceData:
			s.mjpegstream.UpdateJPEG(data)
		case <-s.done:
			return
		}
	}
}

func (s *Stream) stopStream() {
	s.done <- true
}

func NewMjpegStreamer() *MjpegStreamer {
	var streamer MjpegStreamer
	streamer.sources = make(map[string]*Stream)
	streamer.serveMux = http.NewServeMux()

	streamer.serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, indexTemplate)
	})

	return &streamer
}

func newStream(name string, sourceData chan []byte) *Stream {
	var stream Stream
	stream.Name = name
	stream.sourceData = sourceData
	stream.mjpegstream = mjpeg.NewStream()
	return &stream
}

func (s *MjpegStreamer) AddStream(name string, sourceData chan []byte) {
	stream := newStream(name, sourceData)
	s.sources[name] = stream
	go stream.updateStream()
	s.serveMux.Handle(fmt.Sprintf("/%s", name), stream.mjpegstream)
}

func (s *MjpegStreamer) StopStream(name string) {
	stream, ok := s.sources[name]
	if ok {
		stream.stopStream()
		delete(s.sources, name)
	}
}

func (s *MjpegStreamer) StartStream() {
	address := "127.0.0.1:4000"
	log.Printf("Starting data streaming on %s\n", address)
	if err := http.ListenAndServe(address, s.serveMux); nil != err {
		log.Printf("Error: %s", err.Error())
	}
}
