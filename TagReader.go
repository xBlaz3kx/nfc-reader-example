package main

import (
	"encoding/hex"
	"fmt"
	"github.com/clausecker/nfc/v2"
	"github.com/warthog618/gpiod"
	"log"
	"time"
)

type TagReader struct {
	TagChannel       chan string
	reader           *nfc.Device
	ResetPin         int
	DeviceConnection string
}

func (reader *TagReader) init() {
	dev, err := nfc.Open(reader.DeviceConnection)
	if err != nil {
		reader.Reset()
		log.Printf("Cannot communicate with the device: %s \n", err)
		return
	}
	reader.reader = &dev
	err = reader.reader.InitiatorInit()
	if err != nil {
		log.Fatal("Failed to initialize")
		return
	}
}

// Reset Implements the hardware reset by pulling the ResetPin low and then releasing.
func (reader *TagReader) Reset() {
	log.Println("Resetting the reader..")
	//refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	pin, err := c.RequestLine(reader.ResetPin, gpiod.AsOutput(0))
	if err != nil {
		log.Println(err)
		return
	}
	err = pin.SetValue(1)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 400)
	err = pin.SetValue(0)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 400)
	err = pin.SetValue(1)
	time.Sleep(time.Millisecond * 100)
	if err != nil {
		log.Println(err)
		return
	}
}

func (reader *TagReader) Cleanup() {
	defer reader.reader.Close()
}

func (reader *TagReader) ListenForTags() {
	//Initialize the reader
	reader.init()
	var (
		err      error
		tagCount int
		target   nfc.Target
		UID      string
	)
	//Listen for all the modulations specified
	var modulations = []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106},
		{Type: nfc.Felica, BaudRate: nfc.Nbr212},
		{Type: nfc.Felica, BaudRate: nfc.Nbr424},
		{Type: nfc.Jewel, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106},
	}
	for {
		// Poll for 300ms
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
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443bTarget)
				var UIDLen = len(card.ApplicationData)
				var ID = card.ApplicationData
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr212}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr424}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Jewel, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.JewelTarget)
				var ID = card.ID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443biClassTarget)
				var ID = card.UID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			}
			// Send the UID of the tag to main goroutine
			reader.TagChannel <- UID
		}
		time.Sleep(time.Second * 1)
	}
}

func NewTagReader(deviceConnection string, tagChannel chan string, resetPin int) *TagReader {
	return &TagReader{DeviceConnection: deviceConnection, TagChannel: tagChannel, ResetPin: resetPin}
}
