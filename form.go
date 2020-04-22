package main

import (
	"log"

	"github.com/youryharchenko/go-zebkit/zkit"
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
	Win      *zkit.Window
	Root     *zkit.Panel
	Status   *zkit.Panel
	Buttons  *zkit.Panel
	Focus    *zkit.Panel
	ChResult chan FormResult
}

// NewForm -
func NewForm(win *zkit.Window, w int, h int, sizeable bool) (form *zkit.Form) {
	//win := NewWindow(obj)
	win.Layoutable.SetSize(w, h)
	win.SetSizeable(sizeable)

	form = &Form{
		Win:      win,
		Root:     zkit.NewPanel(win.Object().Get("root"), nil),
		Status:   zkit.NewPanel(win.Object().Get("status"), nil),
		Buttons:  zkit.NewPanel(win.Object().Get("buttons"), nil),
		ChResult: make(chan FormResult, 0),
	}
	return
}

// RunModal -
func (form *Form) RunModal(parent *zkit.Layoutable) (res FormResult, err error) {
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
