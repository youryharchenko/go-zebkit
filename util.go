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
