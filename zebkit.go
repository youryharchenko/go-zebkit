package main

import (
	"log"
	"strconv"

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

// MakePassTextField -
func (ui *PkgUI) MakePassTextField(id string, txt string, maxSize int, showLast bool) (pt *TextField) {
	o := ui.Obj.Get("PassTextField").New(txt, maxSize, showLast)
	setID(o, id)
	pt = NewTextField(o)
	return
}

// MakeTextArea -
func (ui *PkgUI) MakeTextArea(id string, text string) (c *TextArea) {
	o := ui.Obj.Get("TextArea").New(text)
	setID(o, id)
	c = NewTextArea(o)
	return
}

// MakeToolBar  -
func (ui *PkgUI) MakeToolBar(id string) (tb *ToolBar) {
	o := ui.Obj.Get("Toolbar").New()
	setID(o, id)
	tb = NewToolBar(o)
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

// MakeDefViews -
func (ui *PkgUI) MakeDefViews() (dv *DefViews) {
	o := ui.Obj.Get("tree").Get("DefViews").New()
	dv = NewDefViews(o)
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

// Kids -
func (ui *PkgUI) Kids(parent *Layoutable) (kids []*Layoutable) {
	objs := parent.Kids()
	for i := 0; i < len(objs); i++ {
		kids = append(kids, NewLayoutable(objs[i]))
	}
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
func (d *PkgData) MakeTreeModel(r *Item) (tm *TreeModel) {
	o := d.Obj.Get("TreeModel").New(r)
	tm = NewTreeModel(o)
	return
}

// MakeItem -
func (d *PkgData) MakeItem(v interface{}) (i *Item) {
	o := d.Obj.Get("Item").New(v)
	i = NewItem(o)
	return
}

// PkgDraw -
type PkgDraw struct {
	Obj *js.Object
}

// NewPkgDraw -
func NewPkgDraw(obj *js.Object) (d *PkgDraw) {
	d = &PkgDraw{
		Obj: obj,
	}
	return d
}

// MakePasswordText -
func (draw *PkgDraw) MakePasswordText(r interface{}) (pt *PasswordText) {
	o := draw.Obj.Get("PasswordText").New(r)
	pt = NewPasswordText(o)
	return
}

// MakeStringRender -
func (draw *PkgDraw) MakeStringRender(name string) (sr *StringRender) {
	o := draw.Obj.Get("StringRender").New(name)
	sr = NewStringRender(o)
	return
}

// Layout -
type Layout interface {
}

// HostDecorativeViews -
type HostDecorativeViews interface {
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

// SetID -
func (l *Layoutable) SetID(id string) {
	setID(l.Obj, id)
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

// Kids -
func (l *Layoutable) Kids() (kids []*js.Object) {
	objs := l.Object().Get("kids")
	if objs == nil || objs.String() == "undefined" {
		return
	}
	n := objs.Length()
	for i := 0; i < n; i++ {
		ind := strconv.Itoa(i)
		kid := objs.Get(ind)
		kids = append(kids, kid)
	}
	return
}

// Remove -
func (l *Layoutable) Remove(obj *js.Object) (r *Layoutable) {
	r = NewLayoutable(l.Object().Call("remove", obj))
	return
}

// Invalidate -
func (l *Layoutable) Invalidate() {
	l.Object().Call("invalidate")
	return
}

// Validate -
func (l *Layoutable) Validate() {
	l.Object().Call("validate")
	return
}

// Fire -
func (l *Layoutable) Fire(eventName string, params interface{}) {
	l.Object().Call("fire", eventName, params)
}

// On -
func (l *Layoutable) On(eventName string, cb interface{}) {
	if len(eventName) == 0 {
		l.Object().Call("on", cb)
	} else {
		l.Object().Call("on", eventName, cb)
	}

}

// Off -
func (l *Layoutable) Off(eventName string, cb interface{}) {
	if len(eventName) == 0 {
		l.Object().Call("off", cb)
	} else {
		l.Object().Call("off", eventName, cb)
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

// SetBackground -
func (p *Panel) SetBackground(v interface{}) {
	p.Object().Call("setBackground", v)
}

// Load -
func (p *Panel) Load(filepath string) {
	p.Object().Call("load", filepath)
}

// SetListLayout -
func (p *Panel) SetListLayout(ax string, gap int) {
	p.Object().Call("setListLayout", ax, gap)
}

// SetFlowLayout -
func (p *Panel) SetFlowLayout(ax string, ay string, dir string, gap int) {
	p.Object().Call("setFlowLayout", ax, ay, dir, gap)
}

// RequestFocus -
func (p *Panel) RequestFocus() {
	p.Object().Call("requestFocus")
}

// Repaint -
func (p *Panel) Repaint() {
	p.Object().Call("repaint")
}

// KeyTyped -
func (p *Panel) KeyTyped(f interface{}) {
	p.Object().Set("keyTyped", f)
}

// ChildKeyTyped -
func (p *Panel) ChildKeyTyped(f interface{}) {
	p.Object().Set("childKeyTyped", f)
}

// ToolBar -
type ToolBar struct {
	Panel
}

// NewToolBar -
func NewToolBar(obj *js.Object) (tb *ToolBar) {
	tb = &ToolBar{}
	tb.Obj = obj
	return
}

// AddImage -
func (tb *ToolBar) AddImage(img interface{}) (p *Panel) {
	p = NewPanel(tb.Object().Call("addImage", img), nil)
	return
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
func (es *EvStatePan) PointerReleased(f interface{}) {
	es.Object().Set("pointerReleased", f)
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

// Fired -
func (b *Button) Fired(f interface{}) {
	b.Object().Set("fired", f)
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

// SetValue -
func (l *Label) SetValue(text string) {
	l.Object().Call("setValue", text)
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

// PassTextField -
type PassTextField struct {
	TextField
}

// NewPassTextField -
func NewPassTextField(obj *js.Object) (pt *PassTextField) {
	pt = &PassTextField{}
	pt.Obj = obj
	return
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
	EventProducer
	HostDecorativeViews
}

// SetModel -
func (bt *BaseTree) SetModel(tm *TreeModel) {
	bt.Object().Call("setModel", tm.Object())
	return
}

// Selected -
func (bt *BaseTree) Selected() (i *Item) {
	i = NewItem(bt.Object().Get("selected"))
	return
}

// Fire -
func (bt *BaseTree) Fire(eventName string, params interface{}) {
	bt.Object().Call("fire", eventName, params)
}

// On -
func (bt *BaseTree) On(eventName string, cb interface{}) {
	bt.Object().Call("on", eventName, cb)
}

// Off -
func (bt *BaseTree) Off(eventName string, cb interface{}) {
	bt.Object().Call("off", eventName, cb)
}

// Select -
func (bt *BaseTree) Select(item *Item) {
	bt.Object().Call("select", item.Object())
}

// SetSelectable -
func (bt *BaseTree) SetSelectable(b bool) {
	bt.Object().Call("setSelectable", b)
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

// SetViewProvider -
func (t *Tree) SetViewProvider(p *DefViews) {
	t.Object().Call("setViewProvider", p)
}

// Tabs -
type Tabs struct {
	Panel
	EventProducer
	HostDecorativeViews
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

// Add -
func (tm *TreeModel) Add(to *Item, an *Item) {
	tm.Object().Call("add", to.Object(), an.Object())
}

// Root -
func (tm *TreeModel) Root() (root *Item) {
	root = NewItem(tm.Object().Get("root"))
	return
}

// FindOne -
func (tm *TreeModel) FindOne(root *Item, value *js.Object) (item *Item) {
	item = NewItem(tm.Object().Call("findOne", root.Object(), value))
	return
}

// Iterate -
func (tm *TreeModel) Iterate(root *Item, f interface{}) {
	tm.Object().Call("iterate", root.Object(), f)
	return
}

// Item -
type Item struct {
	Obj *js.Object
	Val *js.Object
}

// NewItem -
func NewItem(o *js.Object) (i *Item) {
	i = &Item{}
	i.Obj = o
	return
}

// Object -
func (i *Item) Object() (o *js.Object) {
	o = i.Obj
	return
}

// Parent -
func (i *Item) Parent() (p *Item) {
	p = NewItem(i.Object().Get("parent"))
	return
}

// Value -
func (i *Item) Value() (v *js.Object) {
	v = i.Object().Get("value")
	return
}

// Kids -
func (i *Item) Kids() (kids []*js.Object) {
	objs := i.Object().Get("kids")

	if objs == nil || objs.String() == "undefined" {
		return
	}
	n := objs.Length()
	for i := 0; i < n; i++ {
		ind := strconv.Itoa(i)
		kid := objs.Get(ind)
		kids = append(kids, kid)
	}
	return
}

// EventProducer -
type EventProducer interface {
}

// View -
type View struct {
	Obj *js.Object
}

// Object -
func (v *View) Object() (o *js.Object) {
	o = v.Obj
	return
}

// Render -
type Render struct {
	View
}

// BaseTextRender -
type BaseTextRender struct {
	Render
}

// TextRender -
type TextRender struct {
	BaseTextRender
}

// PasswordText -
type PasswordText struct {
	TextRender
}

// NewPasswordText -
func NewPasswordText(obj *js.Object) (pt *PasswordText) {
	pt = &PasswordText{}
	pt.Obj = obj
	return
}

// SetEchoChar -
func (pt *PasswordText) SetEchoChar(ch string) {
	pt.Object().Call("setEchoChar", ch)
}

// BaseViewProvider -
type BaseViewProvider struct {
	Obj *js.Object
}

// Object -
func (bp *BaseViewProvider) Object() (o *js.Object) {
	o = bp.Obj
	return
}

// DefViews -
type DefViews struct {
	BaseViewProvider
}

// NewDefViews -
func NewDefViews(obj *js.Object) (dv *DefViews) {
	dv = &DefViews{}
	dv.Obj = obj
	return
}

// SetView -
func (dv *DefViews) SetView(f interface{}) {
	dv.Object().Set("getView", f)
}

// StringRender -
type StringRender struct {
	BaseTextRender
}

// NewStringRender -
func NewStringRender(obj *js.Object) (sr *StringRender) {
	sr = &StringRender{}
	sr.Obj = obj
	return
}

// Event -
type Event struct {
	Obj *js.Object
	Src *js.Object
}

// Object -
func (e *Event) Object() (o *js.Object) {
	o = e.Obj
	return
}

// Source -
func (e *Event) Source() (s *js.Object) {
	s = e.Src
	return
}

// KeyEvent -
type KeyEvent struct {
	Event
}

// NewKeyEvent -
func NewKeyEvent(obj *js.Object) (ke *KeyEvent) {
	ke = &KeyEvent{}
	ke.Obj = obj
	return
}

// Code -
func (ke *KeyEvent) Code() (c string) {
	c = ke.Object().Get("code").String()
	return
}
