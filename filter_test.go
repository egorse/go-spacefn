package main

import (
	"fmt"
	"testing"

	evdev "github.com/egorse/golang-evdev"
	"github.com/stretchr/testify/assert"
)

func TestFilter1(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		desc  string
		inp   []evdev.InputEvent
		out   []evdev.InputEvent
		state int
	}{
		{
			desc:  "Empty call 1",
			inp:   nil,
			out:   []evdev.InputEvent{},
			state: 0,
		},
		{
			desc:  "Empty call 2",
			inp:   []evdev.InputEvent{},
			out:   []evdev.InputEvent{},
			state: 0,
		},
		{
			desc: "Press/release W",
			inp: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 1}, // press W
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 0}, // release W
				{Type: evdev.SYN_REPORT},
			},
			out: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 1}, // press W
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 0}, // release W
				{Type: evdev.SYN_REPORT},
			},
			state: 0,
		},
		{
			desc: "Press/release 'space'",
			inp: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 1}, // press 'space'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 0}, // release 'space'
				{Type: evdev.SYN_REPORT},
			},
			out: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 1}, // press 'space'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 0}, // release 'space'
				{Type: evdev.SYN_REPORT},
			},
			state: 0,
		},
		{
			desc: "Press 'space', press W, release W, release space",
			inp: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 1}, // press 'space'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 1}, // press W
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 0}, // release W
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 0}, // release 'space'
				{Type: evdev.SYN_REPORT},
			},
			out: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_UP, Value: 1}, // press 'up'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_UP, Value: 0}, // release 'up'
				{Type: evdev.SYN_REPORT},
			},
			state: 0,
		},
		{
			desc: "Press 'space', press W, release space, release W",
			inp: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 1}, // press 'space'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 1}, // press W
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_SPACE, Value: 0}, // release 'space'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_W, Value: 0}, // release W
				{Type: evdev.SYN_REPORT},
			},
			out: []evdev.InputEvent{
				{Type: evdev.EV_KEY, Code: evdev.KEY_UP, Value: 1}, // press 'up'
				{Type: evdev.SYN_REPORT},
				{Type: evdev.EV_KEY, Code: evdev.KEY_UP, Value: 0}, // release 'up'
				{Type: evdev.SYN_REPORT},
			},
			state: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stateState := 0
			state, out := filterEvents(stateState, tC.inp)
			if false {
				for _, ev := range out {
					fmt.Printf("out,%v,%v,0x%x,%v\n", evdev.ByEventType[int(ev.Type)][int(ev.Code)], int(ev.Type), int(ev.Code), ev.Value)
				}
				fmt.Println()
			}
			assert.Equal(tC.state, state)
			assert.Equal(tC.out, out)
		})
	}
}
