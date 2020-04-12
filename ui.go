package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object) {
	zUI, zLayout := NewUI(ui), NewLayout(layout)

	root := zUI.MakeCanvas("z", innerWidth, innerHeight).Root()

	mainMenu := zUI.MakeMenuBar("mainMenu", menuList())
	statusBar := zUI.MakeStatusBarPan("statusBar", 6)

	statusText := zUI.MakeLabel("statusText", "Ready")
	statusBar.Add("left", &statusText.Layoutable)

	tree := zUI.MakeTree("tree", map[string]interface{}{
		"value": "Root",
		"kids": []interface{}{
			"Item 1",
			"Item 2",
		},
	}, true)

	textArea := zUI.MakeTextArea("textArea", "A text ... ")

	splitPan := zUI.MakeSplitPan("splitPan", &tree.Panel, &textArea.Panel, "vertical")
	splitPan.SetLeftMinSize(250)
	splitPan.SetRightMinSize(250)
	splitPan.SetGripperLoc(300)
	splitPan.Properties("", map[string]interface{}{
		"padding": 6,
	})

	button := zUI.MakeButton("button", "Clear")
	button.PointerReleased(func(e *js.Object) {
		log.Println("Click!", e)
		NewTextArea(root.ByPath("#textArea", nil)).SetValue("")
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
