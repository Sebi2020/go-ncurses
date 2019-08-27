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
	AttrNormal      Attribute = C.A_NORMAL
	// Hightlighted text
	AttrHighlighted Attribute = C.A_STANDOUT
	// Underlined text
	AttrUnderline   Attribute = C.A_UNDERLINE
	AttrReversed    Attribute = C.A_REVERSE
	// Blinking text
	AttrBlink       Attribute = C.A_BLINK
	// Dimmed text
	AttrDim         Attribute = C.A_DIM
	// Bold text
	AttrBold        Attribute = C.A_BOLD
	// Protected text
	AttrProtect     Attribute = C.A_PROTECT
	// Hidden text
	AttrInvisible   Attribute = C.A_INVIS
	// Alternative Charset
	AttrAltcharset  Attribute = C.A_ALTCHARSET
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