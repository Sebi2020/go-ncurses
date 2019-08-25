// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <curses.h>
	"C"
	"errors"
	"fmt"
)

// A ncurses base color
type Color int

type pairId int

type pairMap map [string]pairId

var numPairs int = 1
var pairs pairMap = make(pairMap)

const (
    ColorBlack Color = iota
    ColorRed
    ColorGreen
    ColorYellow
    ColorBlue
    ColorMagenta
    ColorCyan
    ColorWhite
)

// Changes terminal to color mode. If the terminal does not support colors. StartColor returns an error.
func StartColor() error {
	if !C.has_colors() {
		return errors.New("Terminal does not support colors")
	}
	C.assume_default_colors(-1,-1)
	pairs["std"] = pairId(0)
	C.start_color()
	return nil
}

// Adds a new pair of colors, which can be used to manipulate the color of terminal outputs.
// Choose one color of type Color.
func AddColorPair(name string, fg,bg Color) {
	pairs[name] = pairId(numPairs)
	C.init_pair(C.short(numPairs),C.short(fg),C.short(bg))
	numPairs++
}

// Selects a font foreground/background color pair. Every write uses the new color pair.
//
// Use AddColorPair to create ncurses Color-Pairs. There is one default color-pair which defaults to the terminal color-set for foreground and background.
func SetColor(name string) error {
	val, ok := pairs[name]
	if !ok {
		return fmt.Errorf("Color \"%s\"pair does not exist!",name)
	}
	GetComChannel() <- Command{
		Name: SETCOLOR,
		Scope:GLOBAL,
		Value:val,
	}
	return nil
}

// Changes foreground/background color-pair of the associated window.
func (w *Window) Wbkgd(pairName string) error {
	val,ok := pairs[pairName]
	if !ok {
		return fmt.Errorf("Color pair \"%s\" does not exist!",pairName)
	}
	w.sendCommand(Command{
		Name:WBKGD,
		Window:w,
		Scope:LOCAL,
		Value:val,
	},w.AutoRefresh)
	return nil
}