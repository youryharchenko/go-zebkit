package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/youryharchenko/go-zebkit/zkit"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object, data *js.Object, draw *js.Object) {

	chOk := make(chan bool, 0)

	zUI, zLayout, zData, zDraw := zkit.NewPkgUI(ui), zkit.NewPkgLayout(layout), zkit.NewPkgData(data), zkit.NewPkgDraw(draw)
	zebkit.Canvas = zUI.MakeCanvas("z", innerWidth, innerHeight)
	rootCanvas = zebkit.Canvas.Root()

	formLogin = MakeFormLogin(zUI, zLayout, zData, zDraw)
	toolBar = zUI.MakeToolBar("toolBar")
	statusBar = zUI.MakeStatusBarPan("statusBar", 6)
	treeModel = zData.MakeTreeModel(MakeAppRoot(zData))
	AddMeta(zData, treeModel)
	tree = zUI.MakeTree("tree", treeModel, true)
	tree.SetSelectable(true)

	go testAndLogin(&rootCanvas.Layoutable, chOk)

	go MakeMainUI(zUI, zLayout, zData, zDraw, &rootCanvas.Layoutable, chOk)

}

// MakeFormLogin -
func MakeFormLogin(zUI *zkit.PkgUI, zLayout *zkit.PkgLayout, zData *zkit.PkgData, zDraw *zkit.PkgDraw) (form *Form) {

	form = NewForm(zUI.MakeWindow("formLogin", "Login", nil), 300, 225, false)

	root := form.Root
	status := form.Status
	buttons := form.Buttons

	root.SetListLayout("left", 6)

	log.Printf("MakeFormLogin - user:%v", session.Data.User)

	userLabel := zUI.MakeLabel("userLabel", "User")
	root.Add("center", &userLabel.Layoutable)
	userField := zUI.MakeTextField("userField", session.Data.User)
	userField.SetHint("User name")
	userField.SetPSByRowsCols(2, 20)
	userField.SetTextAlignment("left")
	root.Add("center", &userField.Layoutable)

	form.Focus = &userField.Panel

	pwdLabel := zUI.MakeLabel("pwdLabel", "Password")
	root.Add("center", &pwdLabel.Layoutable)
	pwdField := zUI.MakePassTextField("pwdField", session.Data.Secret, 10, false)
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

	buttonOK := zUI.MakeButton("buttonOK", "OK")
	root.Add("center", &buttonOK.Layoutable)
	buttonOK.Fired(fnOk)

	status.Remove(status.Kids()[0])
	status.SetFlowLayout("center", "center", "horizontal", 4)

	statLabel := zUI.MakeLabel("statLabel", "Ready")
	status.Add("center", &statLabel.Layoutable)

	close := zkit.NewButton(buttons.Kids()[0])
	close.Fired(fnClose)

	form.ChildKeyTyped(func(key *js.Object) {
		k := zkit.NewKeyEvent(key)
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
func MakeMainUI(zUI *zkit.PkgUI, zLayout *zkit.PkgLayout, zData *zkit.PkgData, zDraw *zkit.PkgDraw, root *zkit.Layoutable, ch chan bool) {
	b := <-ch

	if b {

		err := RefreshTree(zUI, zData)
		if err != nil {
			log.Printf("Refresh tree error: %v", err)
			return
		}

		fillToolBar(toolBar)

		statusBar.SetBackground("lightgrey")
		statusText := zUI.MakeLabel("statusText", "Ready")
		statusBar.Add("left", &statusText.Layoutable)

		defViews := zUI.MakeDefViews()
		defViews.SetView(func(t *js.Object, i *js.Object) (v *js.Object) {
			name := zkit.NewItem(i).Value().Get("name").String()
			v = zDraw.MakeStringRender(name).Object()
			return
		})
		tree.SetViewProvider(defViews)
		tree.On("selected", func(src *js.Object, i *js.Object) {
			//log.Println("Tree node cur selected:", src)
			s := zkit.NewTree(src).Selected()
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
			ks := zkit.NewPanel(src, nil).Kids()
			for i := range ks {
				log.Println(i, ks[i].Get("id").String(), ks[i].Get("state").String())
				if ks[i].Get("state").String() == "over" {
					DispatchToolBarEvent(zUI, zData, ks[i].Get("id").String())
				}
			}
		})

		//textArea1 := ui.MakeTextArea("textArea1", "A text1 ... ")
		//textArea2 := ui.MakeTextArea("textArea2", "A text2 ... ")
		tabs = zUI.MakeTabs("tabs", "top")
		//tabs.Add("Text1", &textArea1.Layoutable)
		//tabs.Add("Text2", &textArea2.Layoutable)

		scrollTreePan := zUI.MakeScrollPan("treeScrollPan", &tree.Panel, "vertical", true)

		splitPan := zUI.MakeSplitPan("splitPan", &scrollTreePan.Panel, &tabs.Panel, "vertical")
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
func DispatchToolBarEvent(zUI *zkit.PkgUI, zData *zkit.PkgData, id string) {
	switch id {
	case "tbRefresh":
		go RefreshTree(zUI, zData)
	case "tbRun":
		go DispatchQuery(zUI, zData)
	}
}

// DispatchQuery -
func DispatchQuery(zUI *zkit.PkgUI, zData *zkit.PkgData) {
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
		srcMap, ok := qry.(map[string]interface{})["source"].(map[string]interface{})
		if !ok {
			log.Println("DispatchQuery error:", errors.New("source is missing"))
			return
		}
		viewMap, ok := qry.(map[string]interface{})["view"].(map[string]interface{})
		if !ok {
			log.Println("DispatchQuery error:", errors.New("view is missing"))
			return
		}

		buf, err := json.Marshal(srcMap)
		if err != nil {
			log.Println("DispatchQuery error:", err)
			return
		}
		src := string(buf)

		log.Println("runQuery:", name, t)
		response, err := runQuery(src)
		if err != nil {
			log.Println("runQuery error:", err)
			return
		}
		log.Printf("Request:\n%s,\nResponse:\n%s", src, response)

		respMap := map[string]interface{}{}
		err = json.Unmarshal([]byte(response), &respMap)
		if err != nil {
			log.Println("DispatchQuery error:", err)
			return
		}
		MakeQryTabGrid(zUI, name, srcMap, viewMap, respMap)
	default:
		return
	}
}

// MakeQryTabGrid -response
func MakeQryTabGrid(zUI *zkit.PkgUI, name string, src map[string]interface{}, view map[string]interface{}, resp map[string]interface{}) {

	result, ok := resp["result"].([]interface{})
	if !ok {
		log.Println("result is missing")
		return
	}

	titles, ok := view["columns"].([]interface{})
	if !ok {
		log.Println("view.columns is missing")
		return
	}

	a := []interface{}{}
	for _, r := range result {
		rMap := r.(map[string]interface{})
		row := []interface{}{}
		for _, t := range titles {
			c, ok := rMap[t.(string)]
			if !ok {
				c = "undefined"
			}
			row = append(row, c)
		}
		a = append(a, row)
	}

	obj := rootCanvas.ByPath("#"+name, nil)
	if obj == nil {
		grid := zUI.MakeGrid(name, a)
		gridStrPan := zUI.MakeGridStretchPan(name+"StretchPan", grid)
		gridCaption := zUI.MakeGridCaption(name+"Caption", titles)
		grid.Add("top", &gridCaption.Layoutable)
		scrollPan := zUI.MakeScrollPan(name+"ScrollPan", &gridStrPan.Panel, "vertical", true)
		tabs.Add(name, &scrollPan.Layoutable)
	} else {
		zkit.NewGrid(obj).SetModel(a)
	}
}

// RefreshTree -
func RefreshTree(zUI *zkit.PkgUI, zData *zkit.PkgData) (err error) {
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

	treeModel = zData.MakeTreeModel(MakeAppRoot(zData))
	AddMeta(zData, treeModel)
	AddObjects(zData, treeModel)
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
func MakeAppRoot(zData *zkit.PkgData) (r *zkit.Item) {
	r = zData.MakeItem(
		js.M{"name": "Application", "type": "app"},
	)
	return
}

// AddMeta -
func AddMeta(zData *zkit.PkgData, model *zkit.TreeModel) (err error) {
	root := model.Root()
	meta := zData.MakeItem(
		js.M{"name": "Meta", "type": "section"},
	)
	model.Add(root, meta)

	traits := zData.MakeItem(
		js.M{"name": "Traits", "type": "section"},
	)
	relations := zData.MakeItem(
		js.M{"name": "Relations", "type": "section"},
	)
	model.Add(meta, traits)
	model.Add(meta, relations)

	body, ok := session.Data.Meta["body"].(map[string]interface{})
	if !ok {
		return
	}

	AddTraits(zData, model, traits, body)
	AddRelations(zData, model, relations, body)
	return
}

// AddObjects -
func AddObjects(zData *zkit.PkgData, model *zkit.TreeModel) (err error) {
	root := model.Root()
	objects := zData.MakeItem(
		js.M{"name": "Objects", "type": "section"},
	)
	model.Add(root, objects)

	a := []interface{}{}
	for _, q := range local.Data.ObjQueries {
		i := map[string]interface{}{"name": q.Name, "source": q.Source, "view": q.View}
		a = append(a, i)
	}
	a = sortItems(a)
	for _, v := range a {
		name, ok := v.(map[string]interface{})["name"]
		if !ok {
			name = "Noname"
		}
		item := zData.MakeItem(
			js.M{"name": name, "type": "obj-query", "data": v},
		)
		model.Add(objects, item)
	}

	return
}

// AddTraits -
func AddTraits(zData *zkit.PkgData, model *zkit.TreeModel, traits *zkit.Item, body map[string]interface{}) (err error) {
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
		item := zData.MakeItem(
			js.M{"name": name, "type": "trait", "data": v},
		)
		model.Add(traits, item)
	}
	return
}

// AddRelations -
func AddRelations(zdata *zkit.PkgData, model *zkit.TreeModel, relations *zkit.Item, body map[string]interface{}) (err error) {
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
		item := zdata.MakeItem(
			js.M{"name": name, "type": "relation", "data": v},
		)
		model.Add(relations, item)
	}
	return
}
