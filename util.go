// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses
import (
	"runtime"
	// #include <curses.h>
	"C"
)

func trace(levelup uint8) (string,string) {
	pc := make([]uintptr, 10)
	if(2+levelup > 10) {
		panic("levelup must be smaller or equal 8")
	}
	n := runtime.Callers(1+int(levelup),pc)
	if n == 0 {
		return "",""
	}
	pc = pc[:n]
	frames :=  runtime.CallersFrames(pc)
	callee,_ := frames.Next()
	caller,_ := frames.Next()
	return caller.Function,callee.Function
}