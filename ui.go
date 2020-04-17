package main

import (
	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object, data *js.Object, draw *js.Object) {

	chOk := make(chan bool, 0)

	zUI, zLayout, zData, zDraw := NewPkgUI(ui), NewPkgLayout(layout), NewPkgData(data), NewPkgDraw(draw)
	root := zUI.MakeCanvas("z", innerWidth, innerHeight).Root()

	//inspectObject("root", root.Object(), 0, 1)

	go func() {
		err := testToken()
		if err != nil {
			formLogin := zUI.MakeFormLogin(zLayout, zData, zDraw)
			i := 0
			for {
				i++
				log.Printf("Login retry %v", i)
				if i > 3 {
					log.Printf("Used max retries: %v", i)
					chOk <- false
					return
				}
				log.Printf("Login - Will RunModal")
				res, err := formLogin.RunModal(&root.Layoutable)
				if err != nil {
					log.Printf("Login - RunModal error:%v", err)
					formLogin.SetValue("#statLabel", fmt.Sprintf("Login error: %v", err))
					continue
				}
				log.Printf("Login - RunModal result:%v", res)
				if res == FormOk {
					err := login()
					if err != nil {
						log.Printf("Login error: %v", err)
						formLogin.SetValue("#statLabel", fmt.Sprintf("Login error: %v", err))
						continue
					}
					break
				}
			}
		}
		log.Printf("Token: %v", session.Data.Token)

		local.Data.User = session.Data.User
		local.Data.Secret = session.Data.Secret
		local.Save()

		chOk <- true
	}()

	go func() {
		b := <-chOk

		if b {
			//mainMenu := zUI.MakeMenuBar("mainMenu", menuList())
			toolBar := zUI.MakeToolBar("toolBar")

			statusBar := zUI.MakeStatusBarPan("statusBar", 6)
			statusBar.SetBackground("lightgrey")
			statusText := zUI.MakeLabel("statusText", "Ready")
			statusBar.Add("left", &statusText.Layoutable)

			treeModel := zData.MakeTreeModel(makeTreeModel())
			tree := zUI.MakeTree("tree", treeModel, true)

			textArea1 := zUI.MakeTextArea("textArea1", "A text1 ... ")
			textArea2 := zUI.MakeTextArea("textArea2", "A text2 ... ")
			tabs := zUI.MakeTabs("tabs", "top")
			tabs.Add("Text1", &textArea1.Layoutable)
			tabs.Add("Text2", &textArea2.Layoutable)

			splitPan := zUI.MakeSplitPan("splitPan", &tree.Panel, &tabs.Panel, "vertical")
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
			root.Properties("", js.M{
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
		}
	}()

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

	//log.Printf("status.Kids len:%v", len(status.Kids()))
	status.Remove(status.Kids()[0])
	status.SetFlowLayout("center", "center", "horizontal", 4)
	//status.SetBackground("grey")
	//log.Printf("status.Kids len:%v", len(status.Kids()))

	statLabel := ui.MakeLabel("statLabel", "Ready")
	//statLabel.SetBackground("white")
	status.Add("center", &statLabel.Layoutable)

	//log.Printf("status.Kids len:%v", len(status.Kids()))

	//inspectObject("status.Kids.0", status.Kids()[0], 0, 1)
	//statLabel := NewLabel(status.Kids()[0])
	//statLabel.SetID("statLabel")

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
