package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// Zebkit -
type Zebkit struct {
	Obj    *js.Object
	Canvas *Canvas
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

// GetPropertyValue -
func (zk *Zebkit) GetPropertyValue(obj *js.Object, path string, useGetter bool) (o *js.Object) {
	o = zk.Obj.Call("getPropertyValue", obj, path, useGetter)
	return
}

// IsAtomic -
func (zk *Zebkit) IsAtomic(obj *js.Object) bool {
	return zk.Obj.Call("isAtomic", obj).Bool()
}

// PkgUI -
type PkgUI struct {
	Obj *js.Object
}

// NewPkgUI -
func NewPkgUI(obj *js.Object) (ui *PkgUI) {
	ui = &PkgUI{
		Obj: obj,
	}
	return ui
}

// MakeCanvas -
func (ui *PkgUI) MakeCanvas(name string, w int, h int) (c *Canvas) {
	o := ui.Obj.Get("zCanvas")
	if len(name) == 0 {
		c = NewCanvas(o.New(w, h))
	} else {
		c = NewCanvas(o.New(name, w, h))
	}
	zebkit.Canvas = c
	return
}

// MakeMenu -
func (ui *PkgUI) MakeMenu(id string, list []interface{}) (m *Menu) {
	o := ui.Obj.Get("Menu").New(list)
	setID(o, id)
	m = NewMenu(o)
	return
}

// MakeMenuBar -
func (ui *PkgUI) MakeMenuBar(id string, list []interface{}) (mb *MenuBar) {
	o := ui.Obj.Get("Menubar").New(list)
	setID(o, id)
	mb = NewMenuBar(o)
	return
}

// MakeStatusBarPan -
func (ui *PkgUI) MakeStatusBarPan(id string, gap int) (sb *StatusBarPan) {
	o := ui.Obj.Get("StatusBarPan").New(gap)
	setID(o, id)
	sb = NewStatusBarPan(o)
	return
}

// MakeMenuItem -
func (ui *PkgUI) MakeMenuItem(item map[string]interface{}) (mi *MenuItem) {
	o := ui.Obj.Get("MenuItem").New(item)
	mi = NewMenuItem(o)
	return
}

// MakeButton -
func (ui *PkgUI) MakeButton(id string, text string) (c *Button) {
	o := ui.Obj.Get("Button").New(text)
	setID(o, id)
	c = NewButton(o)
	return
}

// MakeLabel -
func (ui *PkgUI) MakeLabel(id string, r interface{}) (l *Label) {
	o := ui.Obj.Get("Label").New(r)
	setID(o, id)
	l = NewLabel(o)
	return
}

// MakeTextField -
func (ui *PkgUI) MakeTextField(id string, r interface{}) (tf *TextField) {
	o := ui.Obj.Get("TextField").New(r)
	setID(o, id)
	tf = NewTextField(o)
	return
}

// MakeTextArea -
func (ui *PkgUI) MakeTextArea(id string, text string) (c *TextArea) {
	o := ui.Obj.Get("TextArea").New(text)
	setID(o, id)
	c = NewTextArea(o)
	return
}

// MakeSplitPan -
func (ui *PkgUI) MakeSplitPan(id string, first *Panel, second *Panel, orient string) (c *SplitPan) {
	o := ui.Obj.Get("SplitPan").New(first.Object(), second.Object(), orient)
	setID(o, id)
	c = NewSplitPan(o)
	return
}

// MakeTree -
func (ui *PkgUI) MakeTree(id string, model *TreeModel, b bool) (t *Tree) {
	o := ui.Obj.Get("tree").Get("Tree").New(model.Object(), b)
	setID(o, id)
	t = NewTree(o)
	return
}

// MakeTabs -
func (ui *PkgUI) MakeTabs(id string, orient string) (t *Tabs) {
	o := ui.Obj.Get("Tabs").New(orient)
	setID(o, id)
	t = NewTabs(o)
	return
}

// MakeWindow -
func (ui *PkgUI) MakeWindow(id string, title string, content *Panel) (w *Window) {
	log.Printf("MakeWindow start")
	var o *js.Object
	if content != nil {
		o = ui.Obj.Get("Window").New(title, content.Object())
	} else {
		o = ui.Obj.Get("Window").New(title, nil)
	}
	setID(o, id)
	w = NewWindow(o)
	log.Printf("MakeWindow finish")
	return
}

// MakeFormLogin -
func (ui *PkgUI) MakeFormLogin(zLayout *PkgLayout, zData *PkgData) (form *Form) {

	form = NewForm(ui.MakeWindow("formLogin", "Login", nil), 300, 225, false)

	root := form.Root
	status := form.Status
	buttons := form.Buttons

	root.SetListLayout("left", 6)

	log.Printf("MakeFormLogin - user:%v", session.Data.User)

	userLabel := ui.MakeLabel("userLabel", "User")
	root.Add("center", &userLabel.Layoutable)
	userField := ui.MakeTextField("userField", session.Data.User)
	userField.SetHint("User name")
	userField.SetPSByRowsCols(2, 20)
	userField.SetTextAlignment("left")
	root.Add("center", &userField.Layoutable)

	pwdLabel := ui.MakeLabel("pwdLabel", "Password")
	root.Add("center", &pwdLabel.Layoutable)
	pwdField := ui.MakeTextField("pwdField", session.Data.Secret)
	pwdField.SetHint("Password")
	pwdField.SetPSByRowsCols(2, 20)
	pwdField.SetTextAlignment("left")
	root.Add("center", &pwdField.Layoutable)

	buttonOK := ui.MakeButton("buttonOK", "OK")
	root.Add("center", &buttonOK.Layoutable)
	buttonOK.PointerReleased(func(e *js.Object) {
		log.Println("FormLogin button OK - pointerReleased")
		log.Println("FormLogin will close")
		form.Close()
		session.Data.User = userField.GetValue().String()
		session.Data.Secret = pwdField.GetValue().String()
		session.Save()
		form.ChResult <- FormOk
	})

	statLabel := ui.MakeLabel("statLabel", "Ready")
	status.Insert(0, "left", &statLabel.Layoutable)

	close := NewButton(buttons.Object().Get("kids").Get("0"))
	close.PointerReleased(func(e *js.Object) {
		log.Println("FormLogin button close - pointerReleased")
		log.Println("FormLogin will close")
		form.Close()
		form.ChResult <- FormClose
	})

	form.SetStatus("Input name, passsword and click OK")

	return
}

// PkgLayout -
type PkgLayout struct {
	Obj *js.Object
}

// NewPkgLayout -
func NewPkgLayout(obj *js.Object) (lo *PkgLayout) {
	lo = &PkgLayout{
		Obj: obj,
	}
	return lo
}

// MakeBorderLayout -
func (lo *PkgLayout) MakeBorderLayout(hgap int, vgap int) (l *BorderLayout) {
	o := lo.Obj.Get("BorderLayout").New(hgap, vgap)
	l = NewBorderLayout(o)
	return
}

// PkgData -
type PkgData struct {
	Obj *js.Object
}

// NewPkgData -
func NewPkgData(obj *js.Object) (d *PkgData) {
	d = &PkgData{
		Obj: obj,
	}
	return d
}

// MakeTreeModel -
func (d *PkgData) MakeTreeModel(arg interface{}) (tm *TreeModel) {
	o := d.Obj.Get("TreeModel").New(arg)
	tm = NewTreeModel(o)
	return
}

// Layout -
type Layout interface {
}

// BorderLayout -
type BorderLayout struct {
	Layout
	Obj *js.Object
}

// NewBorderLayout -
func NewBorderLayout(obj *js.Object) (l *BorderLayout) {
	l = &BorderLayout{}
	l.Obj = obj
	return
}

// Object -
func (bl *BorderLayout) Object() (o *js.Object) {
	o = bl.Obj
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
	r = NewPanel(c.Obj.Get("root"), nil)
	return
}

// PathSearch -
type PathSearch interface {
	ByPath(path string, arg interface{})
}

// Layoutable -
type Layoutable struct {
	PathSearch
	EventProducer
	Layout *js.Object
	Obj    *js.Object
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

// Resized -
func (l *Layoutable) Resized(pw int, ph int) {
	l.Object().Call("resized", pw, ph)
}

// SetSize -
func (l *Layoutable) SetSize(w int, h int) {
	l.Object().Call("setSize", w, h)
}

// Add -
func (l *Layoutable) Add(constr interface{}, d *Layoutable) {
	log.Println("Add:", constr)
	l.Object().Call("add", constr, d.Object())
}

// Insert -
func (l *Layoutable) Insert(i int, constr interface{}, d *Layoutable) {
	log.Println("Insert:", constr)
	l.Object().Call("insert", i, constr, d.Object())
}

// SetByConstraints -
func (l *Layoutable) SetByConstraints(constr interface{}, d *Layoutable) {
	log.Println("SetByConstraints:", constr)
	l.Object().Call("setByConstraints", constr, d.Object())
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
func NewPanel(obj *js.Object, layout *js.Object) (p *Panel) {
	p = &Panel{}
	p.Obj = obj
	p.Layout = layout
	return
}

// Load -
func (p *Panel) Load(filepath string) {
	p.Obj.Call("load", filepath)
}

// SetListLayout -
func (p *Panel) SetListLayout(ax string, gap int) {
	p.Object().Call("setListLayout", ax, gap)
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

// GetValue -
func (l *Label) GetValue() (v *js.Object) {
	v = l.Object().Call("getValue")
	return
}

// TextField -
type TextField struct {
	Label
}

// NewTextField -
func NewTextField(obj *js.Object) (tf *TextField) {
	tf = &TextField{}
	tf.Obj = obj
	return
}

// SetValue -
func (tf *TextField) SetValue(text string) {
	tf.Object().Call("setValue", text)
}

// SetHint -
func (tf *TextField) SetHint(text string) {
	tf.Object().Call("setHint", text)
}

// SetPSByRowsCols -
func (tf *TextField) SetPSByRowsCols(r int, c int) {
	tf.Object().Call("setPSByRowsCols", r, c)
}

// SetTextAlignment -
func (tf *TextField) SetTextAlignment(ax string) {
	tf.Object().Call("setTextAlignment", ax)
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

// Tabs -
type Tabs struct {
	Panel
	EventProducer
}

// NewTabs -
func NewTabs(obj *js.Object) (t *Tabs) {
	t = &Tabs{}
	t.Obj = obj
	return
}

// Window -
type Window struct {
	Panel
}

// NewWindow -
func NewWindow(obj *js.Object) (w *Window) {
	w = &Window{}
	w.Obj = obj
	return
}

// SetSizeable -
func (w *Window) SetSizeable(sizeable bool) {
	w.Object().Call("setSizeable", sizeable)
}

// Close -
func (w *Window) Close() {
	w.Object().Call("close")
}

// DataModel -
type DataModel interface {
	Object() *js.Object
}

// TreeModel -
type TreeModel struct {
	DataModel
	EventProducer
	Obj *js.Object
}

// NewTreeModel -
func NewTreeModel(obj *js.Object) (tm *TreeModel) {
	tm = &TreeModel{}
	tm.Obj = obj
	return
}

// Object -
func (tm *TreeModel) Object() (o *js.Object) {
	o = tm.Obj
	return
}

// EventProducer -
type EventProducer interface {
}
