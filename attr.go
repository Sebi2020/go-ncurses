// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <curses.h>
	"C"
	// "fmt"
)
type Attribute int

const (
	// Normal text
	AttrNormal Attribute = C.A_NORMAL
	// Hightlighted text
	AttrHighlighted = C.A_STANDOUT
	// Underlined text
	AttrUnderline = C.A_UNDERLINE
	AttrReversed = C.A_REVERSE
	// Blinking text
	AttrBlink = C.A_BLINK
	// Dimmed text
	AttrDim = C.A_DIM
	// Bold text
	AttrBold = C.A_BOLD
	// Protected text
	AttrProtect = C.A_PROTECT
	// Hidden text
	AttrInvisible = C.A_INVIS
	// Alternative Charset
	AttrAltcharset = C.A_ALTCHARSET
)

// Sets char attributes for following output in Window w.
func (w *Window) SetAttribute(att... Attribute) {
	attSet := 0
	for _,v := range att {
		attSet |= int(v)
	}
	com := Command{
		Name: ATTRSET,
		Scope: LOCAL,
		Window: w,
		Value: attSet,
	}
	GetComChannel() <- com
}