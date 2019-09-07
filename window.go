// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
    // #cgo LDFLAGS: -lncursesw
    // #include <binding.h>
    // #include <ncurses.h>
    "C"
    "fmt"
    "unsafe"
    "errors"
)

// Type which holds information about the current cursor position
type Position struct {
    X uint16
    Y uint16
}

// Alias for position which holds information about terminal size
type Size Position

type winst *C.struct__win_st

func (p1 Position) greater(p2 Position) bool {
    return p1.Y > p2.Y|| p1.X > p2.X
}

// Holds information about the currently used Terminal and works as a handle for all ncurses
// related function calls.
type Window struct {
    name string
    chandle unsafe.Pointer
    max Size
    begin Position
    Cursor Position
    // Set to true, if you want to automatically refresh the window after it recieved a command
    AutoRefresh bool
    AutoCursor bool
    AutoEcho bool
    lastColor string
    inputBuffer []byte
    // Controls maximum character count for reads. This controls the n parameter of ncurses wgetnstr function.
    //
    // See: http://manpages.org/getstr/3
    IBufSize int
    panelList map[string]Panel
}

// Creates a new window. Make sure the windows do not overlap.
//    Notice: Fields for position and Size are swapped. The first parameter sets the column count, not the line count.
// TODO(sebi2020): Add DelWindow method
func NewWindow(name string, begin Position, end Size) (*Window, error) {
    if !initialized {
        return nil,errors.New("Ncurses is not initialized")
    }
    if begin.greater(Position(wins["stdscr"].max)) || Position(end).greater(Position(wins["stdscr"].max)) {
        return nil,errors.New(fmt.Sprintf("Window dimensions greater than terminal size\nMain Size: %v, Window start: %v,Window end:%v",wins["stdscr"].max,begin,end))
    }
    w := &Window{
        name:name,
        begin:begin,
        max:end,
        Cursor:Position{0,0},
        AutoRefresh:false,
        AutoEcho:false,
        lastColor: "std",
        AutoCursor:false,
        IBufSize:255,
    }
    w.chandle = unsafe.Pointer(C.newwin(C.int(begin.Y) + C.int(end.Y),C.int(begin.X) + C.int(end.X),C.int(begin.Y),C.int(begin.X)))
    w.GetMaxYX()
    wins.append(w)
    return w,nil
}

// Implementation of the Stringer Interface for type Window
func (w *Window) String() string { // Implements interface Stringer {
    return fmt.Sprintf("%T{%s %s}", w,w.max,w.Cursor)
}

// Implementation of the Stringer Interface for type Position
func (p Position) String() string {
    return fmt.Sprintf("(%d, %d)", p.Y,p.X);
}

// Implementation of the Stringer Interface for type Size
func (s Size) String() string {
    return fmt.Sprintf("<%d, %d>", s.Y,s.X)
}

func (w *Window) sendCommand(c Command, refresh bool) {
    GetComChannel() <-c
    if refresh {
        w.Refresh()
    }
}

// Retrieves the terminal height and width (gathered through Terminfo & Termcap DB)
func(w *Window) GetMaxYX() (Size,error) {
    if !initialized {
        return Size{},errors.New("Not initialized");
    }
    lines,cols := C.int(0),C.int(0)
    C.bind_get_maxyx((*C.struct__win_st)(w.chandle),&lines,&cols)
    lines_go, cols_go := uint16(lines),uint16(cols)
    w.max  = Size{cols_go,lines_go}
    return w.max,nil
}

// Returns name of the window from which this method is called.
func (w *Window) GetName() string {
    return w.name;
}

// Retrieves one Rune from user.
func (w *Window) Getch() rune {
    if w.AutoEcho {
        SetEcho(true)
    }
    if w.AutoCursor {
        SetCursor(CURSOR_VISIBLE)
    }
    com := Command{
        Name: GETCH,
        Scope: LOCAL,
        Window: w,
        Value: make(chan rune),
    }
    w.sendCommand(com,false)
    ret := <-com.Value.(chan rune)
    if w.AutoEcho {
        SetEcho(false)
    }
    if w.AutoCursor {
        SetCursor(CURSOR_HIDDEN)
    }
    return ret
}

// Moves the the cursor relative to the beginning of w *Window.
func (w *Window) Move(y,x uint16) {
    com := Command{
        Name: MOVE,
        Window:w,
        Value: Position{x,y},
        Scope: LOCAL,
    }
    w.Cursor = Position{x,y}
    w.sendCommand(com,w.AutoRefresh)
}

// Refreshes terminal screen. Outputs all content written since last call to Refresh()
func (w *Window) Refresh() {
    com := Command{
        Name: REFRESH,
        Window:w,
        Scope:LOCAL,
    }
    w.sendCommand(com,false)
}

// Implements the Writer interface.
// Allows you to write strings with fmt.Fprintf(window, format,...args) to the associated window.
func (w *Window) Write(p []byte) (n int, err error) {
    if !initialized {
        return 0,errors.New("ncurses is not initialized")
    }
    n = len(p)
    com := Command{
        Name: ADD,
        Window:w,
        Value:string(p),
        Scope:LOCAL,
    }
    w.sendCommand(com,w.AutoRefresh)
    return n,nil
}

// Inserts a string at the current position.
func (w *Window) Insert(format string, val ...interface{}) {
    com := Command{
        Name: INSERT,
        Window: w,
        Value: fmt.Sprintf(format,val...),
        Scope: LOCAL,
    }
    w.sendCommand(com,w.AutoRefresh)
}

// Allow you to enable window scrolling. Must be called before calling Scroll(n int)
func (w *Window) SetScrolling(enable bool) {
    com := Command{
        Name: SCROLLOK,
        Window: w,
        Value: enable,
        Scope: LOCAL,
    }
    w.sendCommand(com,false)
}
func (w *Window) Scroll(n int) {
    com := Command{
        Name: SCROLL,
        Window: w,
        Value: n,
        Scope: LOCAL,
    }
    w.sendCommand(com,w.AutoRefresh)
}


// Deletes all content of w *Window.
func (w *Window) Clear() {
    com := Command{
        Name: CLEAR,
        Window: w,
        Scope: LOCAL,
    }
    w.sendCommand(com,w.AutoRefresh)
}

func (w *Window) Erase() {
    com := Command{
        Name: ERASE,
        Window: w,
        Scope: LOCAL,
    }
    w.sendCommand(com,w.AutoRefresh)
}

func (w *Window) GetLastColumn() uint16 {
    return w.max.X-1
}

func(w *Window) GetLastLine() uint16 {
    return w.max.Y-1
}

func (w *Window) DrawBG(bg string,p Position, s Size) {
    bx,by := p.X,p.Y
    ex,ey := p.X+s.X,p.Y+s.Y
    colorBefore := w.GetLastColor()
    w.SetColor(bg)
    for Y := by; Y < ey; Y++ {
        str := make([]byte,0,255)
        for X := bx; X < ex; X++ {
            str = append(str,' ')
        }
        w.Move(Y,bx)
        w.Write(str)
    }
    w.SetColor(colorBefore)
}