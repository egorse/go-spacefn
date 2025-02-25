package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	evdev "github.com/egorse/golang-evdev"
	"github.com/egorse/uinput"
)

func handleInputEvents(ch chan inputEvents) {
	// Create virtual keyboard
	path := "/dev/uinput"
	name := "go-spacefn"
	keyboard, err := uinput.CreateKeyboard(path, []byte(name))
	if err != nil {
		log.Fatalf("cannot create virtual keyboard %s, %s: %v", path, name, err)
	}
	defer keyboard.Close()

	// the state represents the fn state
	state := 0 // 0 - normal state, 1 - fn pressed but no decision yet, 2 - fn bypass by repetition, 3 - remap mode

	for ie := range ch {
		now := time.Now().UnixNano()
		events := ie.events

		if monitor {
			for _, ev := range events {
				fmt.Printf("%v,inp,%v,%v,0x%x,%v\n", now, evdev.ByEventType[int(ev.Type)][int(ev.Code)], int(ev.Type), int(ev.Code), ev.Value)
			}
			fmt.Println()
		}

		state, events = filterEvents(state, events)

		if monitor {
			for _, ev := range events {
				fmt.Printf("%v,out,%v,%v,0x%x,%v\n", now, evdev.ByEventType[int(ev.Type)][int(ev.Code)], int(ev.Type), int(ev.Code), ev.Value)
			}
			fmt.Println()
		}

		// bypass
		b := new(bytes.Buffer)
		if err := binary.Write(b, binary.LittleEndian, &events); err != nil {
			log.Fatal(err)
		}
		buffer := b.Bytes()
		_, err := keyboard.File().Write(buffer)
		if err != nil {
			log.Fatal(err)
		}
	}
}
