# go-ncurses
**go-ncurses** is a wrapper for [go](https://www.golang.org) of the famous ncurses library.

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
  w2,err := ncurses.NewWindow("dialog",ncurses.Position{20,5},ncurses.Size{25,5})

  // This can fail if the terminal is too small.
  if err != nil {
    panic(err)
  }
  w2.AutoRefresh = true

  // Use color pair wb (2)
  w2.Wbkgd("wb")

  // Draw a border around our "Greeting Window".
  w2.Box()

  // Move cursor relative to the window borders of w2
  w2.Move(2,3)

  // Output our greeting text
  fmt.Fprintf(w2, "Hello from Go\u2122-Lang!") 

  // Move cursor relative to the beginning of our main window
  w.Move(17,19)

  // Output exit instruction for the user
  fmt.Fprintf(w," => Press a key to exit <=")

  // Wait for user input (e.g. keypress)
  w.Getch()
}
```