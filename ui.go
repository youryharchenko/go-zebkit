package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// BuildUI -
func BuildUI(ui *js.Object, layout *js.Object) {
	zUI, zLayout := NewUI(ui), NewLayout(layout)

	root := zUI.MakeCanvas("", 1280, 800).Root()

	mainMenu := zUI.MakeMenuBar("mainMenu", menuList())

	textArea := zUI.MakeTextArea("textArea", "A text ... ")

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
			"top":    mainMenu.Object(),
			"center": textArea.Object(),
			"bottom": button.Object(),
		},
	})

}
