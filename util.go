package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	jsb "github.com/gopherjs/jsbuiltin"
)

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

func inspectObject(pref string, obj *js.Object, deep int, bound int) {
	if deep > bound {
		return
	}
	for _, v := range js.Keys(obj) {
		p := obj.Get(v)
		np := pref + "." + v
		t := jsb.TypeOf(p)
		switch t {
		case "string", "number", "boolean":
			log.Printf("%s: %v(%s)", np, p, t)
		case "function":
			log.Printf("%s: %v(%s)", np, "()", t)
		case "object":
			log.Printf("%s: %v(%s)", np, "{}", t)
			inspectObject(np, p, deep+1, bound)
		default:
			log.Printf("%s: %v(%s)", np, "unknown", t)
		}

	}
}
