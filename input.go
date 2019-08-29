// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
    "io"
)

// Controls maximum character count for reads. This controls the n parameter of ncurses wgetnstr function.
//
// See: http://manpages.org/getstr/3http://manpages.org/getstr/3
var iBufSize int = 255

type readTuple struct {
	ret chan []byte
	n int
}

// Implements the reader interface
func (w *Window) Read(p []byte) (int,error) {
	if (w.inputBuffer == nil) {
		res := &readTuple{
			ret: make(chan []byte),
		}
		com := Command{
			Name: READSTR,
			Window: w,
			Value: res,
		}
		w.sendCommand(com,false)
		w.inputBuffer = (<-res.ret)[:res.n]
	}
	n := copy(p,w.inputBuffer)
	if(n > 0) {
		w.inputBuffer = w.inputBuffer[n:]
		return n,nil
	} else {
		w.inputBuffer = nil
		return n,io.EOF
	}
}