package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object, data *js.Object) {

	zUI, zLayout, zData := NewPkgUI(ui), NewPkgLayout(layout), NewPkgData(data)

	root := zUI.MakeCanvas("z", innerWidth, innerHeight).Root()
	//inspectObject("root", root.Object(), 0, 1)

	mainMenu := zUI.MakeMenuBar("mainMenu", menuList())
	statusBar := zUI.MakeStatusBarPan("statusBar", 6)

	statusText := zUI.MakeLabel("statusText", "Ready")
	statusBar.Add("left", &statusText.Layoutable)

	treeModel := zData.MakeTreeModel(map[string]interface{}{
		"value": "Root",
		"kids": []interface{}{
			"Item 1",
			"Item 2",
		},
	})
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
	splitPan.Properties("", map[string]interface{}{
		"padding": 6,
	})

	button := zUI.MakeButton("button", "Clear")
	button.PointerReleased(func(e *js.Object) {
		log.Println("Click!", e)
		NewTextArea(root.ByPath("#textArea1", nil)).SetValue("")
	})

	root.Properties("", map[string]interface{}{
		"border":  "plain",
		"padding": 8,
		"layout":  zLayout.MakeBorderLayout(6, 0).Object(),
		"kids": map[string]interface{}{
			"right":  button.Object(),
			"top":    mainMenu.Object(),
			"center": splitPan.Object(),
			"bottom": statusBar.Object(),
		},
	})

}
