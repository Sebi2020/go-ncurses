// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

/*
  Provides an API binding for the TUI-Library libncurses.
  The binding is designed to be used in multithreaded environments and therefore uses
  an command channel to communicate with the ncurses library.
 */
package ncurses

import (
    // #cgo LDFLAGS: -lncursesw
    // #include <binding.h>
    // #include <curses.h>
    // #include <locale.h>
    "C"
    "errors"
)

// checks if ncurses is initialized
var initialized bool = false

type winList map[string]*Window

// Keeps track of allocated windows
var wins winList = make(winList)

func (wl winList) append(w *Window) {
    wl[w.name] = w
}

func (wl winList) delete(w *Window) {
    delete(wl,w.name)
}

// Initializes the ncurses library and install cleanup functions. The function returns a new Term struct.
// After calling this method no outputs should be done with fmt.Print*. Make sure to call Term.Endwin() before making any output with fmt.Print*.
func Initscr() (*Window,error) {
    if(initialized) {
        return nil,errors.New("Already initialized!")
    }
    C.bind_set_locale()
    C.initscr()
    w := &Window{
        name:       "stdscr",
        chandle:    C.bind_get_stdscr(),
        begin:      Position{0,0},
        AutoRefresh: false,
        AutoCursor: false,
        IBufSize: 255,
    }
    C.keypad((*C.struct__win_st)(w.chandle),true)
    C.move(0,0)
    initialized = true
    w.GetMaxYX()
    wins.append(w)
    go processCommands()
    return w,nil;
}

// Closes an initialized terminal. Must be called before the program is about to exit.
func Endwin() error {
    if !initialized {
        return errors.New("Not initialized")
    }

    for _,w := range(wins) {
        wins.delete(w)
    }
    close(GetComChannel())  
    C.endwin();
    initialized = false
    return nil;
}

// If enabled, typed runes are printed to terminal.
func SetEcho(on bool) error {
    if !initialized {
        return errors.New("Not initialized")
    }
    if on {
        C.echo()
    } else {
        C.noecho()
    }
    return nil
}

// Defines visibility of the terminal cursor.
type CursorVisibility uint8

const (
    // Cursor is hidden
    CURSOR_HIDDEN CursorVisibility = iota
    // Cursor is visible
    CURSOR_VISIBLE
    // Cursor is visible and highlighted (not supported on all terminals)
    CURSOR_HIGHTLIGHTED
)

// Sets visiblity of the    terminal cursor.
func SetCursor(choice CursorVisibility) error {
    if !initialized {
        return errors.New("Not initialized")
    }
    C.curs_set(C.int(choice))
    return nil
}