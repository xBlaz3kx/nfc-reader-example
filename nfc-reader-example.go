package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	rfidChannel := make(chan string)
	quitChannel := make(chan os.Signal, 1)
	// Create an abstraction of the Reader, DeviceConnection string is empty if you want the library to autodetect your reader
	rfidReader := NewTagReader("", rfidChannel, 19)
	// Listen for an RFID/NFC tag in another goroutine
	go rfidReader.ListenForTags()
	for {
		fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
		select {
		case tagId := <-rfidReader.TagChannel:
			fmt.Println(tagId)
			continue
		case <-quitChannel:
			rfidReader.Cleanup()
			break
		default:
			time.Sleep(time.Millisecond * 300)
			continue
		}
	}
}
