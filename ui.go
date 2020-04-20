package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object, data *js.Object, draw *js.Object) {

	chOk := make(chan bool, 0)

	zUI, zLayout, zData, zDraw := NewPkgUI(ui), NewPkgLayout(layout), NewPkgData(data), NewPkgDraw(draw)
	rootCanvas = zUI.MakeCanvas("z", innerWidth, innerHeight).Root()

	formLogin = zUI.MakeFormLogin(zLayout, zData, zDraw)
	toolBar = zUI.MakeToolBar("toolBar")
	statusBar = zUI.MakeStatusBarPan("statusBar", 6)
	treeModel = zData.MakeTreeModel(zData.MakeAppRoot())
	zData.AddMeta(treeModel)
	tree = zUI.MakeTree("tree", treeModel, true)
	tree.SetSelectable(true)

	go testAndLogin(&rootCanvas.Layoutable, chOk)

	go zUI.MakeMainUI(zLayout, zData, zDraw, &rootCanvas.Layoutable, chOk)

}

// MakeFormLogin -
func (ui *PkgUI) MakeFormLogin(zLayout *PkgLayout, zData *PkgData, zDraw *PkgDraw) (form *Form) {

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

	form.Focus = &userField.Panel

	pwdLabel := ui.MakeLabel("pwdLabel", "Password")
	root.Add("center", &pwdLabel.Layoutable)
	pwdField := ui.MakePassTextField("pwdField", session.Data.Secret, 10, false)
	pwdField.SetHint("Password")
	pwdField.SetPSByRowsCols(2, 20)
	pwdField.SetTextAlignment("left")
	root.Add("center", &pwdField.Layoutable)

	fnOk := func(e *js.Object) {
		log.Println("FormLogin - OK")
		form.Close()
		session.Data.User = userField.GetValue().String()
		session.Data.Secret = pwdField.GetValue().String()
		session.Save()
		form.ChResult <- FormOk
	}

	fnClose := func(e *js.Object) {
		log.Println("FormLogin - Close")
		form.Close()
		form.ChResult <- FormClose
	}

	buttonOK := ui.MakeButton("buttonOK", "OK")
	root.Add("center", &buttonOK.Layoutable)
	buttonOK.Fired(fnOk)

	status.Remove(status.Kids()[0])
	status.SetFlowLayout("center", "center", "horizontal", 4)

	statLabel := ui.MakeLabel("statLabel", "Ready")
	status.Add("center", &statLabel.Layoutable)

	close := NewButton(buttons.Kids()[0])
	close.Fired(fnClose)

	form.ChildKeyTyped(func(key *js.Object) {
		k := NewKeyEvent(key)
		log.Printf("KeyTyped:%v", k.Code())
		switch k.Code() {
		case "Enter":
			fnOk(key)
		}
	})

	form.SetValue("#statLabel", "Input name, passsword and click OK")

	return
}

// MakeMainUI -
func (ui *PkgUI) MakeMainUI(zLayout *PkgLayout, zData *PkgData, zDraw *PkgDraw, root *Layoutable, ch chan bool) {
	b := <-ch

	if b {

		err := ui.RefreshTree(zData)
		if err != nil {
			log.Printf("Refresh tree error: %v", err)
			return
		}

		fillToolBar(toolBar)

		statusBar.SetBackground("lightgrey")
		statusText := ui.MakeLabel("statusText", "Ready")
		statusBar.Add("left", &statusText.Layoutable)

		defViews := ui.MakeDefViews()
		defViews.SetView(func(t *js.Object, i *js.Object) (v *js.Object) {
			name := NewItem(i).Value().Get("name").String()
			v = zDraw.MakeStringRender(name).Object()
			return
		})
		tree.SetViewProvider(defViews)
		tree.On("selected", func(src *js.Object, i *js.Object) {
			//log.Println("Tree node cur selected:", src)
			s := NewTree(src).Selected()
			if s.Object() != nil {
				v := s.Value()
				name := v.Get("name").String()
				t := v.Get("type").String()
				log.Println("Tree node cur selected:", name, t)
				session.Data.Item = name + " " + t
				session.Save()
			}
		})

		toolBar.On("./*", func(src *js.Object, arg1 *js.Object, arg2 *js.Object) {
			ks := NewPanel(src, nil).Kids()
			for i := range ks {
				log.Println(i, ks[i].Get("id").String(), ks[i].Get("state").String())
				if ks[i].Get("state").String() == "over" {
					ui.DispatchToolBarEvent(zData, ks[i].Get("id").String())
				}
			}
		})

		//textArea1 := ui.MakeTextArea("textArea1", "A text1 ... ")
		//textArea2 := ui.MakeTextArea("textArea2", "A text2 ... ")
		tabs = ui.MakeTabs("tabs", "top")
		//tabs.Add("Text1", &textArea1.Layoutable)
		//tabs.Add("Text2", &textArea2.Layoutable)

		splitPan := ui.MakeSplitPan("splitPan", &tree.Panel, &tabs.Panel, "vertical")
		splitPan.SetLeftMinSize(250)
		splitPan.SetRightMinSize(250)
		splitPan.SetGripperLoc(300)
		splitPan.Properties("", js.M{
			"padding": 6,
		})
		/*
			button := zUI.MakeButton("button", "Clear")
			button.PointerReleased(func(e *js.Object) {
				log.Println("Click!", e)
				NewTextArea(root.ByPath("#textArea1", nil)).SetValue("")
			})
		*/
		rootCanvas.Properties("", js.M{
			"border":  "plain",
			"padding": 8,
			"layout":  zLayout.MakeBorderLayout(6, 0).Object(),
			"kids": js.M{
				//"right":  button.Object(),
				"top":    toolBar.Object(),
				"center": splitPan.Object(),
				"bottom": statusBar.Object(),
			},
		})

		tree.RequestFocus()
	}
}

// DispatchToolBarEvent -
func (ui *PkgUI) DispatchToolBarEvent(zData *PkgData, id string) {
	switch id {
	case "tbRefresh":
		go ui.RefreshTree(zData)
	case "tbRun":
		go ui.DispatchQuery(zData)
	}
}

// DispatchQuery -
func (ui *PkgUI) DispatchQuery(zData *PkgData) {
	item := findItemInTreeModel(treeModel, session.Data.Item)
	if item == nil {
		return
	}
	v := item.Value()
	if v == nil {
		return
	}
	name, t, qry := v.Get("name").String(), v.Get("type").String(), v.Get("data").Interface()
	switch t {
	case "obj-query":
		src, ok := qry.(map[string]interface{})["src"].(string)
		if !ok {
			return
		}
		log.Println("runQuery:", name, t)
		response, err := runQuery(src)
		if err != nil {
			log.Println("runQuery error:", err)
		}
		log.Printf("Request:\n%s,\nResponse:\n%s", src, response)
		ui.MakeQryTabGrid(name, src, response)
	default:
		return
	}
}

// MakeQryTabGrid -response
func (ui *PkgUI) MakeQryTabGrid(name string, src string, response string) {
	respMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(response), &respMap)
	if err != nil {
		log.Println(err)
		return
	}
	result, ok := respMap["result"].([]interface{})
	if !ok {
		log.Println("result is missing")
		return
	}

	a := []interface{}{}
	for _, r := range result {
		rMap := r.(map[string]interface{})
		row := []interface{}{}
		for _, c := range rMap {
			row = append(row, c)
		}
		a = append(a, row)
	}

	obj := rootCanvas.ByPath("#"+name, nil)
	if obj == nil {
		tabs.Add(name, &ui.MakeGrid(name, a).Layoutable)
	} else {
		NewGrid(obj).SetModel(a)
	}
}

// RefreshTree -
func (ui *PkgUI) RefreshTree(zData *PkgData) (err error) {
	r := `{
		"request": {
			"command": "export",
			"service": "meta"
		},
		"db": {
			"driver": "",
			"connection": "",
			"show": true
		},
		"body": {}
	}`

	chOk := make(chan bool, 0)
	go testAndLogin(&rootCanvas.Layoutable, chOk)
	ok := <-chOk
	if !ok {
		return errors.New("test and login failed")
	}

	response, err := post("/api/run", r)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Response:%s", response)

	meta := js.M{}
	err = json.Unmarshal([]byte(response), &meta)
	if err != nil {
		return
	}
	session.Data.Meta = meta
	session.Save()

	treeModel = zData.MakeTreeModel(zData.MakeAppRoot())
	zData.AddMeta(treeModel)
	zData.AddObjects(treeModel)
	tree.SetModel(treeModel)
	//tree.Invalidate()
	if len(session.Data.Item) > 0 {
		tree.Select(findItemInTreeModel(treeModel, session.Data.Item))
		//tree.Repaint()
	}
	//tree.Validate()

	return
}

// MakeAppRoot -
func (data *PkgData) MakeAppRoot() (r *Item) {
	r = data.MakeItem(
		js.M{"name": "Application", "type": "app"},
	)
	return
}

// AddMeta -
func (data *PkgData) AddMeta(model *TreeModel) (err error) {
	root := model.Root()
	meta := data.MakeItem(
		js.M{"name": "Meta", "type": "section"},
	)
	model.Add(root, meta)

	traits := data.MakeItem(
		js.M{"name": "Traits", "type": "section"},
	)
	relations := data.MakeItem(
		js.M{"name": "Relations", "type": "section"},
	)
	model.Add(meta, traits)
	model.Add(meta, relations)

	body, ok := session.Data.Meta["body"].(map[string]interface{})
	if !ok {
		return
	}

	data.AddTraits(model, traits, body)
	data.AddRelations(model, relations, body)
	return
}

// AddObjects -
func (data *PkgData) AddObjects(model *TreeModel) (err error) {
	root := model.Root()
	objects := data.MakeItem(
		js.M{"name": "Objects", "type": "section"},
	)
	model.Add(root, objects)

	a := []interface{}{}
	for _, v := range local.Data.ObjQueries {
		i := map[string]interface{}{"name": v.Name, "src": v.Src}
		a = append(a, i)
	}
	a = sortItems(a)
	for _, v := range a {
		name, ok := v.(map[string]interface{})["name"]
		if !ok {
			name = "Noname"
		}
		item := data.MakeItem(
			js.M{"name": name, "type": "obj-query", "data": v},
		)
		model.Add(objects, item)
	}

	return
}

// AddTraits -
func (data *PkgData) AddTraits(model *TreeModel, traits *Item, body map[string]interface{}) (err error) {
	a, ok := body["traits"].([]interface{})
	if !ok {
		return
	}
	a = sortItems(a)
	for _, v := range a {
		name, ok := v.(map[string]interface{})["name"]
		if !ok {
			name = "Noname"
		}
		item := data.MakeItem(
			js.M{"name": name, "type": "trait", "data": v},
		)
		model.Add(traits, item)
	}
	return
}

// AddRelations -
func (data *PkgData) AddRelations(model *TreeModel, relations *Item, body map[string]interface{}) (err error) {
	a, ok := body["relations"].([]interface{})
	if !ok {
		return
	}
	a = sortItems(a)
	for _, v := range a {
		name, ok := v.(map[string]interface{})["name"]
		if !ok {
			name = "Noname"
		}
		item := data.MakeItem(
			js.M{"name": name, "type": "relation", "data": v},
		)
		model.Add(relations, item)
	}
	return
}
