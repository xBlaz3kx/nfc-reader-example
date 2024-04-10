package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, end := signal.NotifyContext(context.Background(), os.Interrupt)
	defer end()

	// Create an abstraction of the Reader, DeviceConnection string is empty if you want the library to autodetect your reader
	rfidReader := NewTagReader("", 19)
	tagChannel := rfidReader.GetTagChannel()

	// Listen for an RFID/NFC tag in a separate goroutine
	go rfidReader.ListenForTags(ctx)

	for {
		select {
		case tagId := <-tagChannel:
			log.Printf("Read tag: %s \n", tagId)
		case <-ctx.Done():
			err := rfidReader.Cleanup()
			if err != nil {
				log.Fatal("Error cleaning up the reader: ", err.Error())
			}
			break
		default:
			log.Printf("%s: Waiting for a tag \n", time.Now().String())
			time.Sleep(time.Millisecond * 300)
		}
	}
}
