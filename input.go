// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <ncurses.h>
	// #include "binding.h"
	"C"
    "io"
)


type readTuple struct {
	ret chan []byte
	n int
}

// Implements the io.Reader interface
func (w *Window) Read(p []byte) (int,error) {
	if w.AutoCursor {
		SetCursor(CURSOR_VISIBLE)
	}
	if w.AutoEcho {
		SetEcho(true)
	}
	if (w.inputBuffer == nil) {
		res := &readTuple{
			ret: make(chan []byte),
		}
		com := Command{
			Name: READSTR,
			Scope: LOCAL,
			Window: w,
			Value: res,
		}
		w.sendCommand(com,false)
		w.inputBuffer = (<-res.ret)[:res.n]
	}
	
	n := copy(p,w.inputBuffer)
	if w.AutoCursor {
			SetCursor(CURSOR_HIDDEN)
	}
	if w.AutoEcho {
		SetEcho(false)
	}

	if(n > 0) {
		w.inputBuffer = w.inputBuffer[n:]
		return n,nil
	} else {
		w.inputBuffer = nil
		return n,io.EOF
	}
}

type Key rune

const (
	KEY_UP = C.KEY_UP
	KEY_DOWN = C.KEY_DOWN
	KEY_LEFT = C.KEY_LEFT
	KEY_RIGHT = C.KEY_RIGHT
	KEY_END = C.KEY_END
	KEY_HOME = C.KEY_HOME
	KEY_BACKSPACE = C.KEY_BACKSPACE
	KEY_PAGEUP = C.KEY_PPAGE
	KEY_PAGEDOWN = C.KEY_NPAGE
	KEY_ESC = '\x1b'
)

var (
	KEY_F1 rune = 0
	KEY_F2 rune = 0
	KEY_F3 rune = 0
	KEY_F4 rune = 0
	KEY_F5 rune = 0
	KEY_F6 rune = 0
	KEY_F7 rune = 0
	KEY_F8 rune = 0
	KEY_F9 rune = 0
	KEY_F10 rune = 0
	KEY_F11 rune = 0
	KEY_F12 rune = 0
)

func init() {
	KEY_F1 = rune(C.bind_fkey(1))
	KEY_F2 = rune(C.bind_fkey(2))
	KEY_F3 = rune(C.bind_fkey(3))
	KEY_F4 = rune(C.bind_fkey(4))
	KEY_F5 = rune(C.bind_fkey(5))
	KEY_F6 = rune(C.bind_fkey(6))
	KEY_F7 = rune(C.bind_fkey(7))
	KEY_F8 = rune(C.bind_fkey(8))
	KEY_F9 = rune(C.bind_fkey(9))
	KEY_F10 = rune(C.bind_fkey(10))
	KEY_F11 = rune(C.bind_fkey(11))
	KEY_F12 = rune(C.bind_fkey(12))
}