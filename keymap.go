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

	// Fx keys
	evdev.KEY_1: evdev.KEY_F1,
	evdev.KEY_2: evdev.KEY_F2,
	evdev.KEY_3: evdev.KEY_F3,
	evdev.KEY_4: evdev.KEY_F4,
	evdev.KEY_5: evdev.KEY_F5,
	evdev.KEY_6: evdev.KEY_F6,
	evdev.KEY_7: evdev.KEY_F7,
	evdev.KEY_8: evdev.KEY_F8,
	evdev.KEY_9: evdev.KEY_F9,
	evdev.KEY_0: evdev.KEY_F10,

	// Esc and `
	evdev.KEY_ESC:   evdev.KEY_GRAVE,
	evdev.KEY_GRAVE: evdev.KEY_ESC,
}
