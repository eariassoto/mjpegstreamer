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
		<img src="/mjpeg">
	</body>
	</html>
	`
)

var (
	stream     *mjpeg.Stream
	streamChan chan []byte
)

// StreamSource provides the source data. Data will be feeded to the streamer
// through the data channel
type StreamSource interface {
	StartSourceStream(dataChan chan<- []byte)
}

// StartStream initiates data stream in address host:port
func StartStream(host string, port int, source StreamSource) {
	address := fmt.Sprintf("%s:%d", host, port)
	streamChan = make(chan []byte)

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start source data update function
	go source.StartSourceStream(streamChan)

	go updateStream(streamChan)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, indexTemplate)
	})

	serveMux.Handle("/mjpeg", stream)

	log.Printf("Starting data streaming on %s\n", address)
	if err := http.ListenAndServe(address, serveMux); nil != err {
		log.Printf("Error: %s", err.Error())
	}
}

func updateStream(streamChan <-chan []byte) {
	for {
		data := <-streamChan
		stream.UpdateJPEG(data)
	}
}
