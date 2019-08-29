// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

package ncurses

import (
    "io"
)

// Implements the reader interface
func (w *Window) Read(p []byte, n int) {
	readComplete := make(chan struct{})
	// TODO
}