package main

import "github.com/gopherjs/gopherjs/js"

// JSONstringify -
func JSONstringify(obj *js.Object) string {
	return js.Global.Get("JSON").Call("stringify", obj).String()
}

func setID(obj *js.Object, id string) {
	if len(id) > 0 {
		obj.Call("setId", id)
	}
}

func initResize() {
	win := js.Global.Get("window")
	canvas := js.Global.Get("document").Call("getElementById", "z")

	resizeCanvas := func() {
		prevWidth := innerWidth
		prevHeight := innerHeight

		innerWidth = win.Get("innerWidth").Int()
		innerHeight = win.Get("innerHeight").Int()

		canvas.Set("width", innerWidth)
		canvas.Set("height", innerHeight)

		if zebkit != nil && zebkit.Canvas != nil {
			zebkit.Canvas.SetSize(innerWidth, innerHeight)
			zebkit.Canvas.Resized(prevWidth, prevHeight)
		}
	}

	win.Call("addEventListener", "resize", resizeCanvas, false)
	// Draw canvas border for the first time.
	resizeCanvas()
}
