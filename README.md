# go-ncurses
[![GoDoc](https://godoc.org/github.com/Sebi2020/go-ncurses?status.svg)](https://godoc.org/github.com/Sebi2020/go-ncurses)

**go-ncurses** is a wrapper for [go](https://www.golang.org) of the famous ncurses library.

## Features

+ **Thread Safety** (e.g. Go-Routine safety)
  
  It's safe to run operations on different window from different go-routines
+ **Auto-Refresh** on output

  You can enable ***Auto-Refresh***, which takes care about window refreshes after writes.
+ **Auto-Cursor** on input

  Automatically enables the input cursor on input if enabled

+ Formatted output(bold, italic, underlined and reversed)

  go-ncurses implemented a Writer allowing text-formatting:
  **Example**:
  ```go
  w,_ = ncurses.Initscr()
  defer ncurses.Endwin()
  w.AutoRefresh = true
  wf := ncurses.NewFormatWriter(w)
  fmt.Fprintf(wf,"-Hello *World*-!")
  ```
  **Output**: *Hello **World***!

## Example
```go
package ncurses

import ( 
	"fmt"
)

main() {
 w,_ := ncurses.Initscr()
  
    // Ensure, that ncurses will be properly exited
  defer ncurses.Endwin()
  defer func() {
    if r := recover(); r != nil {
      ncurses.Endwin()
      fmt.Printf("panic:\n%s\n", r)
      os.Exit(-1)
    }
  }()

  // Enable color mode
  ncurses.StartColor()

  // Define color pairs
  ncurses.AddColorPair("bw", ncurses.ColorGreen,ncurses.ColorBlack)
  ncurses.AddColorPair("wb",ncurses.ColorWhite, ncurses.ColorBlue)

  // Set cursor visiblity to hidden
  ncurses.SetCursor(ncurses.CURSOR_HIDDEN)

  // Automatically refresh after each command
  w.AutoRefresh = true

  // Set color for stdscr-window to system defaults.
  w.Wbkgd("std")

  // Draw a border around main window (stdscr)
  w.Box()

  // Create a new window for greeting-text at cell (x=20,y=5) with a size of 25 x 5 cells.
  w2,err := ncurses.NewWindow("dialog",ncurses.Position{20,5},ncurses.Size{25,6})

  // This can fail if the terminal is too small.
  if err != nil {
    panic(err)
  }

  w2.AutoRefresh = true

  // Show cursor on input operation
  w2.AutoCursor = true
  w.AutoCursor = true

  // Use color pair wb (2)
  w2.Wbkgd("wb")

  // Draw a border around our "Greeting Window".
  w2.Box()

  // Move cursor relative to the window borders of w2
  w2.Move(2,3)

  // Output our greeting text
  fmt.Fprintf(w2, "Hello from Go™-Lang!") 

  // Create an input field label
  w2.Move(3,4)
  fmt.Fprintf(w2,"Name ❯ ")

  // Create an input field, restricted to 10 chars (ASCII, less for Unicode)
  w2.SetAttribute(ncurses.AttrUnderline)
  fmt.Fprintf(w2,"          ")
  w2.IBufSize = 10
  w2.Move(3,11)
  v := make([]byte,10)

  n,err := w2.Read(v)
  
  // Move cursor relative to the beginning of our main window
  w.Move(12,26)

  // Use go-ncurses formats for output formatting
  wf := ncurses.NewFormatWriter(w)

  // Name will be displayed with bold-italic font.
  fmt.Fprintf(wf,"Hello -*%s*-!",v[:n])

  // Output exit instruction for the user
  w.Move(17,19)
  fmt.Fprintf(w," => Press a key to exit <=")
  
  
  // Wait for user input (e.g. keypress)
  w.Getch()
}
```