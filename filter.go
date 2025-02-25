package main

import (
	evdev "github.com/egorse/golang-evdev"
)

func filterEvents(state int, events []evdev.InputEvent) (int, []evdev.InputEvent) {
	oe := make([]evdev.InputEvent, 0, 16)
	for _, ev := range events {
		t := ev.Time
		key := -1
		if int(ev.Type) == evdev.EV_KEY {
			key = int(ev.Code)
		}
		is_released, is_pressed, is_repeated := ev.Value == 0, ev.Value == 1, ev.Value == 2
		is_syn := ev.Type == evdev.SYN_REPORT
		is_fn_key := key == fnKey
		is_fn_pressed := is_fn_key && is_pressed
		is_fn_released := is_fn_key && is_released
		is_fn_repeated := is_fn_key && is_repeated

		/*
			state == -1
			almost normal
		*/
		if state == -1 {
			state = 0 // enter normal mode
			if is_syn {
				// skip syn report right after release fn key
				continue
			}
		}
		if state != 1 && is_fn_released {
			// The fn key released not in normal state
			// enter almost normal
			state = -1
			continue
		}
		/*
			state == 0
			normal
		*/
		if state == 0 && is_fn_pressed {
			state = 1 // fn pressed but yet no decision
			continue  // skip current event
		}
		if state == 0 {
			// In normal state we still might have some keys be remapped from previous presses
			// Keep those well remapped until released
			// That would ensure we dont have races with fn key
			n, ok := remappedKeys[key]
			if ok {
				ev.Code = uint16(n)
				if is_released {
					delete(remappedKeys, key)
				}
			}
		}

		/*
			state == 1
			fn pressed but yet no decision
		*/
		if state == 1 && is_fn_released {
			oe = append(oe, evdev.InputEvent{Time: t, Type: ev.Type, Code: ev.Code, Value: 1}) // append extra key press
			oe = append(oe, evdev.InputEvent{Time: t, Type: evdev.SYN_REPORT})                 // syn report
			state = 0                                                                          // normal
		}
		if state == 1 && is_fn_repeated {
			oe = append(oe, evdev.InputEvent{Time: t, Type: ev.Type, Code: ev.Code, Value: 1}) // append extra key press
			oe = append(oe, evdev.InputEvent{Time: t, Type: evdev.SYN_REPORT})                 // syn report
			state = 2                                                                          // fn bypass by repetition
		}
		if state == 1 && !is_fn_key && is_pressed {
			n, ok := keyMap[key]
			if ok {
				ev.Code = uint16(n)
				state = 3 // remap mode
			}
		}
		if state == 1 && is_syn {
			// skip syn report as we pending decision and there is no key press in the output event added yet
			continue
		}

		/*
			state == 3
			remap mode
		*/
		if state == 3 && is_fn_key {
			continue // skip current event
		}
		if state == 3 && !is_fn_key {
			// WARN in remap mode we probably should skip non remapped keys but lets keep those for a while
			n, ok := keyMap[key]
			if ok {
				ev.Code = uint16(n)

				// Maintain map of activelly remappend keys
				// Those has to be released in normal state
				// as well repetition shall works well
				if is_pressed || is_repeated {
					remappedKeys[key] = n
				} else if is_released {
					delete(remappedKeys, key)
				}
			}
		}

		oe = append(oe, ev)
	}

	return state, oe
}

var remappedKeys = make(map[int]int)
