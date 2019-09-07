// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
	// #include <ncurses.h>
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
		return fmt.Sprintf("Unkown (0x%x)",int(a))
	}
	return ""
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
	default:
		return fmt.Sprintf("Unkwown (0x%x)",int(st))
	}
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
	Add bool
}

func initOpenStates(m modList) map[rune]bool {
	open := make(map[rune]bool)
	for _,v := range m {
		open[v] = true
	}
	return open
}

func ParseFormatStr(s string) []interface{} {
	STATE := A_STATE_CHAR
	output := make([]rune,0,255)
	input := ([]rune)(s)
	ops := make([]interface{},0,8)
	Iidx := 0
	open := initOpenStates(modifiers)

	for Iidx < len(input) {
		v := input[Iidx]
		switch STATE {
			case A_STATE_CHAR:
				if modifiers.contains(v) {
					STATE = A_STATE_MODIFIER
				} else {
					output = append(output,v)
				}
				Iidx++
			case A_STATE_MODIFIER:
				if input[Iidx-1] == v {
					STATE = A_STATE_ESCAPE
				} else {
					STATE = A_STATE_FORMAT_FLAG
				}
			case A_STATE_FORMAT_FLAG:
				ops = append(ops,string(output))
				output = make([]rune,0,255)
				ops = append(ops,attrOp{
					Attr: modMapping[input[Iidx-1]],
					Add: open[input[Iidx-1]],
				})
				open[input[Iidx-1]] = !open[input[Iidx-1]]
				STATE = A_STATE_CHAR
			case A_STATE_ESCAPE:
				output = append(output,v)
				Iidx++
				STATE = A_STATE_CHAR
		}
	}
	ops = append(ops,string(output))
	return ops
}

func (aw *AttributeWriter) Write([] byte) (int,error) {
	if !initialized {
		return 0,errors.New("ncurses is not initialized!")
	}
	return 0,io.EOF
}

// Allows you to use go-ncurses formats to format output.
//
// Format specifiers
//
// go-ncurses formats support the following specifiers:
//    *text*: Bold font
//    ~text~: Reversed (colored) font
//    -text-: Italic font
//    _text_: Underlined text
//
// Escaping
//
// Double a format specifier to escape it.
//
type FormatWriter Window


// Returns a pointer to a new FormatWriter
func NewFormatWriter(w *Window) *FormatWriter {
	return (*FormatWriter)(w)
}

// Implements the io.Writer interface
func(fw *FormatWriter) Write(p []byte) (int,error) {
	ops := ParseFormatStr(string(p))
	attr := AttrNormal
	cSum := 0
	for _,v := range ops {
		switch t := v.(type) {
		case string:
			n,err := (*Window)(fw).Write([]byte(v.(string)))
			cSum += n
			if err != nil {
				return cSum,err
			}
		case attrOp:
			if t.Add {
				attr |= t.Attr
			} else {
				attr &= ^t.Attr
			}
			(*Window)(fw).SetAttribute(attr)
		}
	}
	return cSum,nil
}