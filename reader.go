package main

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	"github.com/clausecker/nfc/v2"
	"github.com/pkg/errors"
	"github.com/warthog618/gpiod"
)

var modulations = []nfc.Modulation{
	{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106},
	{Type: nfc.Felica, BaudRate: nfc.Nbr212},
	{Type: nfc.Felica, BaudRate: nfc.Nbr424},
	{Type: nfc.Jewel, BaudRate: nfc.Nbr106},
	{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106},
}

type Reader interface {
	ListenForTags(ctx context.Context) error
	Cleanup() error
	Reset()
	GetTagChannel() <-chan string
}

type TagReader struct {
	TagChannel       chan string
	reader           *nfc.Device
	ResetPin         int
	DeviceConnection string
}

func (reader *TagReader) init() error {
	dev, err := nfc.Open(reader.DeviceConnection)
	if err != nil {
		// Reset the reader if there is an error
		reader.Reset()
		return errors.Wrap(err, "Cannot communicate with the device")
	}

	reader.reader = &dev
	err = reader.reader.InitiatorInit()
	if err != nil {
		return errors.Wrap(err, "Cannot initialize the reader")
	}

	return nil
}

func NewTagReader(deviceConnection string, resetPin int) *TagReader {
	return &TagReader{DeviceConnection: deviceConnection, TagChannel: make(chan string, 10), ResetPin: resetPin}
}

// Reset performs a hardware reset by pulling the ResetPin low and then releasing.
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

func (reader *TagReader) Cleanup() error {
	return reader.reader.Close()
}

func (reader *TagReader) GetTagChannel() <-chan string {
	return reader.TagChannel
}

func (reader *TagReader) getIdFromTarget(target nfc.Target) (*string, error) {
	var UID string
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
	default:
		return nil, errors.New("Unknown modulation")
	}

	return &UID, nil
}

func (reader *TagReader) ListenForTags(ctx context.Context) error {
	//Initialize the reader
	err := reader.init()
	if err != nil {
		return errors.Wrap(err, "Cannot initialize the reader")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Poll for 300ms
			tagCount, target, err := reader.reader.InitiatorPollTarget(modulations, 1, 300*time.Millisecond)
			if err != nil {
				log.Println("Error polling the reader", err)
				continue
			}

			// Check if a tag was detected
			if tagCount > 0 {
				// Get the UID of the tag based on the modulation type
				id, err := reader.getIdFromTarget(target)
				if err != nil {
					log.Println("Error getting ID from target", err)
					continue
				}

				// Send the UID of the tag to main goroutine
				reader.TagChannel <- *id
			}

			time.Sleep(time.Second * 1)
		}
	}
}
