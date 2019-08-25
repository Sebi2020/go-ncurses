// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

// TODO(sebi2020) Ensure that only one instance of Term control the Terminal
import (
    // #cgo LDFLAGS: -lncursesw
    // #include <binding.h>
    // #include <curses.h>
    "C"
    "fmt"
    "os"
)

// Command channel ensures Thread safety
var comChan chan Command = make(chan Command,4)

// Name of a ncurses command
type CommandName uint8

// Data of a ncurses command
type CommandValue interface{}

// Scope of Command
type CommandScope int

const (
    GLOBAL CommandScope = iota          // not a window specific command
    LOCAL                               // a window specific command
)

//All available commands, which can be passed to the command channel
const (
    MOVE CommandName = iota // Move the terminal cursor
    ADD                     // Add string at current cursor location
    INSERT                  // Insert string at current cursor location
    DELETE                  // Delete string at current cursor locatio
    REFRESH                 // Flush buffer to video memory
    CLEAR                   // Clear the entire window
    SCROLLOK                // Enables scrolling
    SCROLL                  // Scrolls current window
    SETCOLOR                // Sets the color for text and background
    WBKGD                   // Sets fg and bg of entire window
    START_TA                // UNIMPL(sebi2020) Initiates a transaction
    END_TA                  // UNIMPL(sebi2020) Finalizes a transaction
)

// TODO(sebi2020): Use generator for this
func (cn CommandName) String() string {
    switch cn {
        case MOVE:
            return "MOV"
        case ADD:
            return "ADD"
        case INSERT:
            return "INS"
        case DELETE:
            return "DEL"
        case REFRESH:
            return "REFRESH"
        case CLEAR:
            return "CLEAR"
        case SCROLLOK:
            return "SCROLLOK"
        case SCROLL:
            return "SCROLL"
        case SETCOLOR:
            return "SETCOLOR"
        case WBKGD:
            return "WBKGD"
        default:
            return fmt.Sprintf("Unkown (%x)",int(cn))
    }
}

// Ncurse related commands
type Command struct {
    // Name of the command
    Name CommandName
    // Window on which the command should be executed
    Window *Window
    // Scope of command
    Scope CommandScope
    // Data which should be passed along with the command
    Value CommandValue
}

// Implements Stringer interface for 'type Command'.
func (c Command) String() string {
    return fmt.Sprintf("Command <%s>",c.Name)
}

// Returns the command channel associated with stdscr (see ncurses documentation for information about stdscr)
func GetComChannel() chan<- Command {
    return comChan
}

func (com Command) execute () {
    var handle winst
    if com.Scope != GLOBAL {
        if com.Window == nil {
            panic(fmt.Sprintf("No Context for command %s", com))
        }
        handle = (*C.struct__win_st)(com.Window.chandle)
    }
    // TODO(sebi2020): Maybe use an interface with type inference to distinguish between different commands.
    switch com.Name {
        case MOVE:
            pos := com.Value.(Position)
            C.wmove(handle,C.int(pos.Y),C.int(pos.X))
        case ADD:
            text := C.CString(com.Value.(string))
            C.bind_waddstr(handle,text)
        case INSERT:
            text := C.CString(com.Value.(string))
            C.winsstr(handle,text)
        case REFRESH:
            C.wrefresh(handle)
        case CLEAR:
            C.wclear(handle)
        case SCROLLOK:
            C.scrollok(handle,C.bool(com.Value.(bool)))
        case SCROLL:
            C.wscrl(handle,C.int(com.Value.(int)))
        case SETCOLOR:
            C.bind_color_set(C.short(com.Value.(pairId)))
        case WBKGD:
            C.bind_wbkgd(handle,C.short(com.Value.(pairId)))
        default:
            panic(fmt.Sprintf("Command %s not implemented",com))
    }
}

func recoverFromPanic() {
    if r := recover(); r != nil {
        Endwin()
        fmt.Printf("Ncurses panic: %s",r)
        os.Exit(-1)
    }
}
func processCommands() {
    defer recoverFromPanic()
    for com := range comChan {
        com.execute()
    }
}
