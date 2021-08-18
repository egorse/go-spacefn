package main

import (
	"log"
	"os"
	"sync"
	"time"

	evdev "github.com/egorse/golang-evdev"
)

var device_glob = "/dev/input/by-id/*"
var monitor = false
var fnKey = evdev.KEY_SPACE
var keyMap = keyMap1

func main() {

	// appears we might start up faster than input device flushed
	// fully by the host and thats triggers some odd repetitions
	time.Sleep(1 * time.Second)

	// List all devices matching glob
	devices, err := evdev.ListInputDevices(device_glob)
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		// You might have to run as sudo or being in proper group (i.e. input)
		log.Fatal("no input devices opened!!! permissions issues?")
	}

	ch := make(chan inputEvents, 32)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		handleInputEvents(ch)
	}()

	var wg2 sync.WaitGroup
	for _, dev := range devices {

		if !hasKeyboardCapabilities(dev) {
			dev.File.Close()
			continue
		}
		log.Printf("using %s, %s", dev.Fn, dev.Name)

		wg2.Add(1)
		// Next goroutines will push all input events to channel to be handled
		go func(dev *evdev.InputDevice) {
			defer wg2.Done()

			err := dev.Grab()
			if err != nil {
				log.Fatalf("cannot grab exclusively %s, %s: %v", dev.Fn, dev.Name, err)
			}
			defer dev.Release()

			for {
				events, err := dev.Read()

				// WARN In case of plug and play devices we have to handle it more nicer
				if err != nil {
					log.Fatalf("error read %s, %s: %v", dev.Fn, dev.Name, err)
					os.Exit(1)
				}

				ie := inputEvents{dev, events}
				ch <- ie
			}
		}(dev)
	}
	wg2.Wait()

	close(ch)
	wg1.Wait()
}

// The hasKeyboardCapabilities return true, if the device has keyboard capabilities
// needed for this application.
// This might select mouses as well, so technically we could even remap mouse buttons.
// But for performance reasons we will try to ignore pointing devices (EV_REL, EV_ABS).
// ATM its not clear would it affects some combo keyboards (i.e. notebook?)
func hasKeyboardCapabilities(dev *evdev.InputDevice) bool {
	requires := []int{evdev.EV_MSC, evdev.EV_KEY, evdev.EV_SYN}
	avoids := []int{evdev.EV_REL, evdev.EV_ABS}

	hasCapability := func(cap int) bool {
		for c := range dev.Capabilities {
			if c.Type == cap {
				return true
			}
		}
		return false
	}

	for _, cap := range requires {
		if !hasCapability(cap) {
			return false
		}
	}

	for _, cap := range avoids {
		if hasCapability(cap) {
			return false
		}
	}

	return true
}

type inputEvents struct {
	dev    *evdev.InputDevice
	events []evdev.InputEvent
}
