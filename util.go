package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
	jsb "github.com/gopherjs/jsbuiltin"
	xhr "github.com/rocketlaunchr/gopherjs-xhr"
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

func initSession() {
	user := local.Data.User
	password := local.Data.Secret
	session = OpenSession()
	if session == nil {
		session = NewSession(SessionData{
			User:   user,
			Secret: password,
		})
		log.Printf("New session:%s", user)
		return
	} else if len(session.Data.User) == 0 {
		session.Data.User = user
		session.Data.Secret = password
		log.Printf("Change user/password session:%s", user)
		return
	}
	log.Printf("Open session:%s", user)
}

func initLocal() {
	local = OpenLocal()
	log.Printf("Open local storage")
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

func login() (err error) {
	req := xhr.NewRequest("POST", "/login")
	req.ResponseType = xhr.Text // Returns response as string
	req.SetRequestHeader("Content-Type", xhr.ApplicationJSON)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, session.Data.User, session.Data.Secret)
	err = req.Send(context.Background(), body)
	if err != nil {
		return
	}
	if !req.IsStatus2xx() {
		err = fmt.Errorf("%v, %s", req.Status, req.StatusText)
		return
	}
	resp := req.ResponseBytes()
	sb := js.M{}
	err = json.Unmarshal(resp, &sb)
	if err != nil {
		return
	}
	t, ok := sb["token"].(string)
	if !ok {
		err = fmt.Errorf("token is missing: %s", string(resp))
		return
	}
	session.Data.Token = t
	session.Save()
	return
}

func testToken() (err error) {
	req := xhr.NewRequest("GET", "/auth-test")
	req.ResponseType = xhr.Text // Returns response as string
	req.SetRequestHeader("Content-Type", xhr.ApplicationJSON)
	req.SetRequestHeader("Authorization", "Bearer "+session.Data.Token)

	//body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, session.Data.User, session.Data.Secret)
	err = req.Send(context.Background(), "")
	if err != nil {
		return
	}
	if !req.IsStatus2xx() {
		err = fmt.Errorf("%v, %s", req.Status, req.StatusText)
		return
	}
	resp := req.ResponseBytes()
	log.Println(string(resp))
	return
}
