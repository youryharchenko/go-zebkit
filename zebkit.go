package main

import "github.com/gopherjs/gopherjs/js"

// Zebkit -
type Zebkit struct {
	Obj *js.Object
}

// NewZebkit -
func NewZebkit() (zk *Zebkit) {
	zk = &Zebkit{
		Obj: js.Global.Get("zebkit"),
	}
	return
}

// UIConfig -
func (zk *Zebkit) UIConfig(args ...interface{}) {
	zk.Obj.Get("ui").Call("config", args...)
}

// Require -
func (zk *Zebkit) Require(args ...interface{}) {
	zk.Obj.Call("require", args...)
}

// UI -
type UI struct {
	Obj *js.Object
}

// NewUI -
func NewUI(obj *js.Object) (ui *UI) {
	ui = &UI{
		Obj: obj,
	}
	return ui
}

// MakeCanvas -
func (ui *UI) MakeCanvas(name string, w int, h int) (c *Canvas) {
	o := ui.Obj.Get("zCanvas")
	if len(name) == 0 {
		c = NewCanvas(o.New(w, h))
	} else {
		c = NewCanvas(o.New(name, w, h))
	}
	return
}

// MakeMenu -
func (ui *UI) MakeMenu(id string, list []interface{}) (m *Menu) {
	o := ui.Obj.Get("Menu").New(list)
	setID(o, id)
	m = NewMenu(o)
	return
}

// MakeMenuBar -
func (ui *UI) MakeMenuBar(id string, list []interface{}) (mb *MenuBar) {
	o := ui.Obj.Get("Menubar").New(list)
	setID(o, id)
	mb = NewMenuBar(o)
	return
}

// MakeStatusBarPan -
func (ui *UI) MakeStatusBarPan(id string, gap int) (sb *StatusBarPan) {
	o := ui.Obj.Get("StatusBarPan").New(gap)
	setID(o, id)
	sb = NewStatusBarPan(o)
	return
}

// MakeMenuItem -
func (ui *UI) MakeMenuItem(item map[string]interface{}) (mi *MenuItem) {
	o := ui.Obj.Get("MenuItem").New(item)
	mi = NewMenuItem(o)
	return
}

// MakeButton -
func (ui *UI) MakeButton(id string, text string) (c *Button) {
	o := ui.Obj.Get("Button").New(text)
	setID(o, id)
	c = NewButton(o)
	return
}

// MakeLabel -
func (ui *UI) MakeLabel(id string, r interface{}) (l *Label) {
	o := ui.Obj.Get("Label").New(r)
	setID(o, id)
	l = NewLabel(o)
	return
}

// MakeTextArea -
func (ui *UI) MakeTextArea(id string, text string) (c *TextArea) {
	o := ui.Obj.Get("TextArea").New(text)
	setID(o, id)
	c = NewTextArea(o)
	return
}

// MakeSplitPan -
func (ui *UI) MakeSplitPan(id string, first *Panel, second *Panel, orient string) (c *SplitPan) {
	o := ui.Obj.Get("SplitPan").New(first.Object(), second.Object(), orient)
	setID(o, id)
	c = NewSplitPan(o)
	return
}

// MakeTree -
func (ui *UI) MakeTree(id string, model interface{}, b bool) (t *Tree) {
	o := ui.Obj.Get("tree").Get("Tree").New(model, b)
	setID(o, id)
	t = NewTree(o)
	return
}

// Layout -
type Layout struct {
	Obj *js.Object
}

// NewLayout -
func NewLayout(obj *js.Object) (lo *Layout) {
	lo = &Layout{
		Obj: obj,
	}
	return lo
}

// Object -
func (lo *Layout) Object() (o *js.Object) {
	o = lo.Obj
	return
}

// MakeBorderLayout -
func (lo *Layout) MakeBorderLayout(hgap int, vgap int) (l *BorderLayout) {
	o := lo.Object().Get("BorderLayout").New(hgap, vgap)
	l = NewBorderLayout(o)
	return
}

// BorderLayout -
type BorderLayout struct {
	Layout
}

// NewBorderLayout -
func NewBorderLayout(obj *js.Object) (l *BorderLayout) {
	l = &BorderLayout{}
	l.Obj = obj
	return
}

// Canvas -
type Canvas struct {
	Panel
}

// NewCanvas -
func NewCanvas(obj *js.Object) (c *Canvas) {
	c = &Canvas{}
	c.Obj = obj
	return
}

// Root -
func (c *Canvas) Root() (r *Panel) {
	r = NewPanel(c.Obj.Get("root"))
	return
}

// PathSearch -
type PathSearch interface {
	ByPath(path string, arg interface{})
}

// Layoutable -
type Layoutable struct {
	PathSearch
	Obj *js.Object
}

// NewLayoutable -
func NewLayoutable(obj *js.Object) (p *Layoutable) {
	p = &Layoutable{}
	p.Obj = obj
	return
}

// Object -
func (l *Layoutable) Object() (o *js.Object) {
	o = l.Obj
	return
}

// ByPath -
func (l *Layoutable) ByPath(path string, arg interface{}) (o *js.Object) {
	if arg == nil {
		o = l.Object().Call("byPath", path)
	} else {
		o = l.Object().Call("byPath", path, arg)
	}
	return
}

// Add -
func (l *Layoutable) Add(constr interface{}, d *Layoutable) {
	l.Object().Call("add", constr, d.Object())
}

// Properties -
func (l *Layoutable) Properties(path string, props map[string]interface{}) {
	if len(path) == 0 {
		l.Object().Call("properties", props)
	} else {
		l.Object().Call("properties", path, props)
	}
}

// Panel -
type Panel struct {
	Layoutable
}

// NewPanel -
func NewPanel(obj *js.Object) (p *Panel) {
	p = &Panel{}
	p.Obj = obj
	return
}

// Load -
func (p *Panel) Load(filepath string) {
	p.Obj.Call("load", filepath)
}

// SplitPan -
type SplitPan struct {
	Panel
}

// NewSplitPan -
func NewSplitPan(obj *js.Object) (sp *SplitPan) {
	sp = &SplitPan{}
	sp.Obj = obj
	return
}

// SetLeftMinSize -
func (sp *SplitPan) SetLeftMinSize(m int) {
	sp.Object().Call("setLeftMinSize", m)
}

// SetRightMinSize -
func (sp *SplitPan) SetRightMinSize(m int) {
	sp.Object().Call("setRightMinSize", m)
}

// SetGripperLoc -
func (sp *SplitPan) SetGripperLoc(l int) {
	sp.Object().Call("setGripperLoc", l)
}

// ViewPan -
type ViewPan struct {
	Panel
}

// StatePan -
type StatePan struct {
	ViewPan
}

// TrackInputEventState -
type TrackInputEventState interface {
	PointerReleased(arg interface{})
}

// EvStatePan -
type EvStatePan struct {
	StatePan
	TrackInputEventState
}

// PointerReleased -
func (es *EvStatePan) PointerReleased(arg interface{}) {
	es.Object().Set("pointerReleased", arg)
}

// Button -
type Button struct {
	EvStatePan
}

// NewButton -
func NewButton(obj *js.Object) (b *Button) {
	b = &Button{}
	b.Obj = obj
	return
}

// Label -
type Label struct {
	ViewPan
}

// NewLabel -
func NewLabel(obj *js.Object) (l *Label) {
	l = &Label{}
	l.Obj = obj
	return
}

// TextField -
type TextField struct {
	Label
}

// SetValue -
func (tf *TextField) SetValue(text string) {
	tf.Object().Call("setValue", text)
}

// TextArea -
type TextArea struct {
	TextField
}

// NewTextArea -
func NewTextArea(obj *js.Object) (ta *TextArea) {
	ta = &TextArea{}
	ta.Obj = obj
	return
}

// BaseList -
type BaseList struct {
	Panel
}

// CompList  -
type CompList struct {
	BaseList
}

// Menu -
type Menu struct {
	CompList
}

// NewMenu -
func NewMenu(obj *js.Object) (m *Menu) {
	m = &Menu{}
	m.Obj = obj
	return
}

// MenuBar -
type MenuBar struct {
	Menu
}

// NewMenuBar -
func NewMenuBar(obj *js.Object) (mb *MenuBar) {
	mb = &MenuBar{}
	mb.Obj = obj
	return
}

// MenuItem -
type MenuItem struct {
	Panel
}

// NewMenuItem -
func NewMenuItem(obj *js.Object) (mi *MenuItem) {
	mi = &MenuItem{}
	mi.Obj = obj
	return
}

// StatusBarPan -
type StatusBarPan struct {
	Panel
}

// NewStatusBarPan -
func NewStatusBarPan(obj *js.Object) (sb *StatusBarPan) {
	sb = &StatusBarPan{}
	sb.Obj = obj
	return
}

// BaseTree -
type BaseTree struct {
	Panel
}

// Tree -
type Tree struct {
	BaseTree
}

// NewTree -
func NewTree(obj *js.Object) (t *Tree) {
	t = &Tree{}
	t.Obj = obj
	return
}
