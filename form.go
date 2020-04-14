package main

import "github.com/gopherjs/gopherjs/js"

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
	ChResult chan FormResult
}

// NewForm -
func NewForm(obj *js.Object) (form *Form) {
	win := NewWindow(obj)
	form = &Form{
		Win:      win,
		ChResult: make(chan FormResult, 0),
	}
	return
}

// RunModal -
func (form *Form) RunModal(parent *Layoutable) (res FormResult, err error) {
	parent.SetByConstraints("center", &form.Win.Layoutable)

	//inspectObject("win.buttons.kids", form.Win.Object().Get("buttons").Get("kids"), 0, 2)

	res = <-form.ChResult
	return
}

// SetStatus -
func (form *Form) SetStatus(stat string) {
	form.Win.ByPath("#statLabel", nil).Call("setValue", stat)
}
