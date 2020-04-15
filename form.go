package main

import (
	"log"
)

// FormResult -
type FormResult int

//
const (
	FormOk FormResult = iota
	FormCancel
	FormClose
)

// Form -
type Form struct {
	Win      *Window
	Root     *Panel
	Status   *Panel
	Buttons  *Panel
	ChResult chan FormResult
}

// NewForm -
func NewForm(win *Window, w int, h int, sizeable bool) (form *Form) {
	//win := NewWindow(obj)
	win.Layoutable.SetSize(w, h)
	win.SetSizeable(sizeable)

	form = &Form{
		Win:      win,
		Root:     NewPanel(win.Object().Get("root"), nil),
		Status:   NewPanel(win.Object().Get("status"), nil),
		Buttons:  NewPanel(win.Object().Get("buttons"), nil),
		ChResult: make(chan FormResult, 0),
	}
	return
}

// RunModal -
func (form *Form) RunModal(parent *Layoutable) (res FormResult, err error) {
	log.Printf("Login - RunModal started")

	parent.SetByConstraints("center", &form.Win.Layoutable)
	res = <-form.ChResult

	log.Printf("Login - RunModal finish")
	return
}

// SetStatus -
func (form *Form) SetStatus(stat string) {
	form.Win.ByPath("#statLabel", nil).Call("setValue", stat)
}

func (form *Form) Close() {
	form.Win.Close()
}
