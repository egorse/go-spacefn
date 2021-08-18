package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"syscall"

	evdev "github.com/egorse/golang-evdev"
	"github.com/egorse/uinput"
)

func trace(str string) {
	if !monitor {
		return
	}
	fmt.Println(str)
}

func handleInputEvents(ch chan inputEvents) {
	// Create virtual keyboard
	path := "/dev/uinput"
	name := "go-spacefn"
	keyboard, err := uinput.CreateKeyboard(path, []byte(name))
	if err != nil {
		log.Fatalf("cannot create virtual keyboard %s, %s: %v", path, name, err)
	}
	defer keyboard.Close()

	state := 0 // 0 - normal state, 1 - fn pressed but no decision yet, 2 - fn bypass by repetition, 3 - remap mode

	for ie := range ch {
		events := ie.events

		if monitor {
			fmt.Println("input:")
			for _, ev := range events {
				fmt.Printf("%v %v(%v, 0x%x) %v\n", ev.Time, evdev.ByEventType[int(ev.Type)][int(ev.Code)], int(ev.Type), int(ev.Code), ev.Value)
			}
			fmt.Println()
		}

		oe := make([]evdev.InputEvent, 0, 16)
		for _, ev := range events {
			ev.Time = syscall.Timeval{Sec: 0, Usec: 0}

			key := -1
			if int(ev.Type) == evdev.EV_KEY {
				key = int(ev.Code)
			}

			switch {
			case state != 1 && key == fnKey && ev.Value == 0 /* released */ :
				{
					state = 0 // normal
				}
			case state == 0 /* normal */ && key == fnKey && ev.Value == 1 /* pressed */ :
				{
					state = 1 // fn pressed but yet no decision
					continue  // skip current event
				}
			case state == 1 /* fn pressed but yet no decision */ && key == fnKey && ev.Value == 0 /* released */ :
				{
					oe = append(oe, evdev.InputEvent{Type: ev.Type, Code: ev.Code, Value: 1}) // append extra key press
					oe = append(oe, evdev.InputEvent{Type: evdev.SYN_REPORT})                 // syn report
					state = 0                                                                 // normal
				}
			case state == 1 /* fn pressed but yet no decision */ && key == fnKey && ev.Value == 2 /* repeated */ :
				{
					oe = append(oe, evdev.InputEvent{Type: ev.Type, Code: ev.Code, Value: 1}) // append extra key press
					oe = append(oe, evdev.InputEvent{Type: evdev.SYN_REPORT})                 // syn report
					state = 2                                                                 // fn bypass by repetition
				}
			case state == 1 /* fn pressed but yet no decision */ && key != fnKey:
				{
					n, ok := keyMap[key]
					if ok {
						ev.Code = uint16(n)
						state = 3 // remap mode
					}

				}
			case state == 3 /* remap mode */ && key == fnKey:
				{
					continue // skip current event
				}
			case state == 3 /* remap mode */ && key != fnKey:
				{
					// WARN in remap mode we probably should skip non remapped keys but lets keep those for a while
					n, ok := keyMap[key]
					if ok {
						ev.Code = uint16(n)
					}
				}
			}

			oe = append(oe, ev)
		}
		events = oe

		if monitor {
			fmt.Println("bypass:")
			for _, ev := range events {
				fmt.Printf("%v %v(%v, 0x%x) %v\n", ev.Time, evdev.ByEventType[int(ev.Type)][int(ev.Code)], int(ev.Type), int(ev.Code), ev.Value)
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
