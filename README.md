# mpjegstreamer

mpjegstream is a simple Go library to stream images through HTTP. It only requires an address and a function to feed the source data channel.

## Example
This program streams a X11 based desktop:
```go
package main

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/eariassoto/mjpegstreamer"
)

type desktopDataSource struct{}

func (d *desktopDataSource) StartSourceStream(dataChan chan<- []byte) {
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
	host := "127.0.0.1"
	port := 4000
	var d desktopDataSource
	mjpegstreamer.StartStream(host, port, &d)
}
```
To watch the stream point your browser to:
```bash
http://127.0.0.1:4000
```
