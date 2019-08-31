// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <curses.h>
	"C"
	"fmt"
	"io"
	"errors"
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
	AttrItalic      Attribute = C.A_ITALIC
	// Protected text
	AttrProtect     Attribute = C.A_PROTECT
	// Hidden text
	AttrInvisible   Attribute = C.A_INVIS
	// Alternative Charset
	AttrAltcharset  Attribute = C.A_ALTCHARSET
)

func (a Attribute) String() string {
	switch(a) {
	case AttrNormal:
		return "AttrNomal"
	case AttrBold:
		return "AttrBold"
	case AttrItalic:
		return "AttrItalic"
	case AttrUnderline:
		return "AttrUnderline"
	case AttrReversed:
		return "AttrReversed"
	default:
		return fmt.Sprintf("Unkown (%x)",a)
	}
}

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

type AttributeWriter struct{
	w io.Writer
}

type attrParserState int

const (
	A_STATE_CHAR attrParserState = iota
	A_STATE_MODIFIER
	A_STATE_ESCAPE
	A_STATE_FORMAT_FLAG
)

func (st attrParserState) String() string {
	switch(st) {
	case A_STATE_CHAR:
		return "A_STATE_CHAR"
	case A_STATE_MODIFIER:
		return "A_STATE_MODIFIER"
	case A_STATE_ESCAPE:
		return "A_STATE_ESCAPE"
	case A_STATE_FORMAT_FLAG:
		return "A_STATE_FORMAT_FLAG"
	}
	return ""
}
type modList []rune
func (m modList) contains(r rune) bool {
	for _,v := range ([]rune)(m) {
		if(v == r) {
			return true
		}
	}
	return false
}

var modifiers modList = modList{'*','_','-','~'}

var modMapping map[rune] Attribute = map[rune] Attribute{
	'*':AttrBold,
	'~':AttrReversed,
	'-':AttrItalic,
	'_':AttrUnderline,
}

type attrOp struct {
	Attr Attribute
	Offset int
	Add bool
}
func initOpenStates(m modList) map[rune]bool {
	open := make(map[rune]bool)
	for _,v := range m {
		open[v] = true
	}
	return open
}

func ParseFormatStr(s string) (string,[]attrOp) {
	STATE := A_STATE_CHAR
	output := make([]rune,0,255)
	input := ([]rune)(s)
	ops := make([]attrOp,0,10)
	Iidx,Oidx := 0,0
	open := initOpenStates(modifiers)

	for Iidx < len(input) {
		v := input[Iidx]
		switch STATE {
			case A_STATE_CHAR:
				if modifiers.contains(v) {
					STATE = A_STATE_MODIFIER
				} else {
					output = append(output,v)
					Oidx++
				}
				Iidx++
			case A_STATE_MODIFIER:
				if input[Iidx-1] == v {
					STATE = A_STATE_ESCAPE
				} else {
					STATE = A_STATE_FORMAT_FLAG
				}
			case A_STATE_FORMAT_FLAG:
				ops = append(ops,attrOp{
					Attr: modMapping[input[Iidx-1]],
					Offset: Oidx,
					Add: open[input[Iidx-1]],
				})
				open[input[Iidx-1]] = !open[input[Iidx-1]]
				STATE = A_STATE_CHAR
			case A_STATE_ESCAPE:
				output = append(output,v)
				Iidx++
				Oidx++
				STATE = A_STATE_CHAR
		}
	}
	return string(output),ops
}

func (aw *AttributeWriter) Write([] byte) (int,error) {
	if !initialized {
		return 0,errors.New("ncurses is not initialized!")
	}
	return 0,io.EOF
}