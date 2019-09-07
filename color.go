// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <ncurses.h>
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
    ColorBrightBlack
    ColorBrightRed
    ColorBrightGreen
    ColorBrightYellow
    ColorBrightBlue
    ColorBrightMagenta
    ColorBrightCyan
    ColorBrightWhite
)

// Changes terminal to color mode. If the terminal does not support colors. StartColor returns an error.
func StartColor() error {
	if !C.has_colors() {
		return errors.New("Terminal does not support colors")
	}
	C.start_color()
	if C.can_change_color() {
		InitColor("blk",CustomColor{Red:0,Green:0,Blue:0})
	}
	pairs["std"] = pairId(-1)
	return nil
}

// Adds a new pair of colors, which can be used to manipulate the color of terminal outputs.
// Choose one color of type Color.
func InitPair(name string, fg,bg Color) {
	pairs[name] = pairId(numPairs)
	C.init_pair(C.short(numPairs),C.short(fg),C.short(bg))
	numPairs++
}

// Selects a font foreground/background color pair. Every write uses the new color pair.
//
// Use AddColorPair to create ncurses Color-Pairs. There is one default color-pair which defaults to the terminal color-set for foreground and background.

func (w *Window) SetColor(name string) error {
	val, ok := pairs[name]
	if !ok {
		return fmt.Errorf("Color \"%s\"pair does not exist!",name)
	}
	w.lastColor = name
	GetComChannel() <- Command{
		Name: SETCOLOR,
		Scope:LOCAL,
		Window: w,
		Value:val,
	}
	return nil
}

func (w *Window) GetLastColor() string {
	return w.lastColor
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

type CustomColor struct {
	Red uint16
	Green uint16
	Blue uint16
	colorId uint16
}

var numColors uint16 = 16
var customColors map[string]CustomColor = make(map[string]CustomColor)

func InitColor(name string, c CustomColor) error {
	if !initialized {
		return errors.New("ncurses it not initalized")
	}
	if !C.can_change_color() {
		return errors.New("Terminal does not support custom colors")
	}
	if c.Red > 999 || c.Green > 999 || c.Blue > 999 || c.Red < 0 || c.Green < 0 || c.Blue < 0 {
		return errors.New("Color values must be in range of 0-999")
	}
	if C.int(numColors) >= C.COLORS {
		return errors.New("Maximum number of custom colors reached")
	}
	c.colorId = numColors
	C.init_color(C.short(c.colorId), C.short(c.Red), C.short(c.Green),C.short(c.Blue))
	customColors[name] = c
	numColors++
	return nil
}

func GetColor(name string) Color {
	val,ok := customColors[name]
	if !ok {
		panic(fmt.Sprintf("Color '%s' does not exists!",name))
	}
	return Color(val.colorId)
}