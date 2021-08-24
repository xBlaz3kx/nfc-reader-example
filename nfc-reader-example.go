package main

import (
	"encoding/hex"
	"fmt"
	"github.com/clausecker/nfc/v2"
	"log"
	"os"
	"time"
)

type TagReader struct {
	TagChannel       chan string
	reader           *nfc.Device
	DeviceConnection string
}

func (reader *TagReader) init() {
	dev, err := nfc.Open(reader.DeviceConnection)
	if err != nil {
		log.Fatalf("Cannot communicate with the device: %s", err)
		return
	}
	reader.reader = &dev
	err = reader.reader.InitiatorInit()
	if err != nil {
		log.Fatal("Failed to initialize")
		return
	}
}

func (reader *TagReader) Cleanup() {
	defer reader.reader.Close()
}

func (reader *TagReader) ListenForTags() {
	//Initialize the reader
	reader.init()
	//Listen for all the modulations specified
	var (
		err      error
		tagCount int
		target   nfc.Target
		UID      string
	)
	var modulations = []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106},
		{Type: nfc.Felica, BaudRate: nfc.Nbr212},
		{Type: nfc.Felica, BaudRate: nfc.Nbr424},
		{Type: nfc.Jewel, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106},
	}
	for {
		// Poll once for 300ms
		tagCount, target, err = reader.reader.InitiatorPollTarget(modulations, 1, 300*time.Millisecond)
		if err != nil {
			fmt.Println("Error polling the reader", err)
			continue
		}
		// Check if any tag was detected
		if tagCount > 0 {
			fmt.Printf(target.String())
			// Transform the target to a specific tag Type and send the UID to the channel
			switch target.Modulation() {
			case nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443aTarget)
				var UIDLen = card.UIDLen
				var ID = card.UID
				// Transform the UID to string and cut the excess bytes
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			case nfc.Modulation{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443bTarget)
				var UIDLen = len(card.ApplicationData)
				var ID = card.ApplicationData
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr212}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr424}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			case nfc.Modulation{Type: nfc.Jewel, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.JewelTarget)
				var ID = card.ID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			case nfc.Modulation{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443biClassTarget)
				var ID = card.UID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:])
				UID = UID[:UIDLen]
				break
			}
			// Send the UID of the tag to main goroutine
			reader.TagChannel <- UID
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	rfidChannel := make(chan string)
	quitChannel := make(chan os.Signal, 1)
	// Create an abstraction of the Reader, DeviceConnection string is empty if you want the library to autodetect your reader
	rfidReader := &TagReader{TagChannel: rfidChannel, DeviceConnection: ""}
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
