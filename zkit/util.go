package zkit

import "github.com/gopherjs/gopherjs/js"

func setID(obj *js.Object, id string) {
	if len(id) > 0 {
		obj.Call("setId", id)
	}
}
