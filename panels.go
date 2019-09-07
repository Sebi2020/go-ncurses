package ncurses


import (
    // #include <ncurses.h>
    // #include <panel.h>
    // #cgo LDFLAGS: -lpanel
    "C"
)

// A C Panel Reference
type Panel struct {
	handle *C.struct_panel
	w *Window
}

func (w *Window) NewPanel(name string) *Panel {
	p := Panel{w:w}
	p.handle = C.new_panel(winst(w.chandle))
	if w.panelList == nil {
		w.panelList = make(map[string]Panel)
	}
	w.panelList[name] = p
	return &p
}

func (p Panel) Up() {
	com := Command{
		Name:PANELUP,
		Scope:LOCAL,
		Window: p.w,
		Value:p,
	}
	p.w.sendCommand(com,p.w.AutoRefresh)
}
func (p Panel) Down() {
	com := Command{
		Name:PANELDOWN,
		Scope:LOCAL,
		Window: p.w,
		Value:p,
	}
	p.w.sendCommand(com,p.w.AutoRefresh)
}