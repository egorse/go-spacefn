package main

import evdev "github.com/egorse/golang-evdev"

var keyMap1 = map[int]int{
	// arrows
	evdev.KEY_W: evdev.KEY_UP,
	evdev.KEY_A: evdev.KEY_LEFT,
	evdev.KEY_S: evdev.KEY_DOWN,
	evdev.KEY_D: evdev.KEY_RIGHT,

	// pgup/pgdown
	evdev.KEY_R: evdev.KEY_PAGEUP,
	evdev.KEY_F: evdev.KEY_PAGEDOWN,

	// home/end
	evdev.KEY_Q: evdev.KEY_HOME,
	evdev.KEY_E: evdev.KEY_END,
}
