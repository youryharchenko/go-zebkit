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
	Focus    *Panel
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
	form.RequestFocus()
	res = <-form.ChResult

	log.Printf("Login - RunModal finish")
	return
}

// SetValue -
func (form *Form) SetValue(id string, stat interface{}) {
	form.Win.ByPath(id, nil).Call("setValue", stat)
}

// RequestFocus -
func (form *Form) RequestFocus() {
	form.Focus.RequestFocus()
}

// KeyTyped -
func (form *Form) KeyTyped(f interface{}) {
	//form.Win.KeyTyped(f)
	form.Root.KeyTyped(f)
}

// ChildKeyTyped -
func (form *Form) ChildKeyTyped(f interface{}) {
	//form.Win.KeyTyped(f)
	form.Root.ChildKeyTyped(f)
}

// Close -
func (form *Form) Close() {
	form.Win.Close()
}
