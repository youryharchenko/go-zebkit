package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
		log.Printf("Change session user/password :%s", user)
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
			if !strings.Contains(v, "parent") {
				inspectObject(np, p, deep+1, bound)
			}
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

func post(uri string, r string) (response string, err error) {
	req := xhr.NewRequest("POST", uri)
	req.ResponseType = xhr.Text // Returns response as string
	req.SetRequestHeader("Content-Type", xhr.ApplicationJSON)
	req.SetRequestHeader("Authorization", "Bearer "+session.Data.Token)

	err = req.Send(context.Background(), r)
	if err != nil {
		return
	}
	if !req.IsStatus2xx() {
		err = fmt.Errorf("%v, %s", req.Status, req.StatusText)
		return
	}
	response = req.ResponseText
	return
}

func refreshMeta() (err error) {
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
	return
}

func makeTreeModel() (tm js.M) {

	tm = js.M{
		"value": "Application",
		"kids": js.S{
			js.M{
				"value": "Meta",
				"kids":  js.S{},
			},
		},
	}

	body, ok := session.Data.Meta["body"].(map[string]interface{})
	if !ok {
		return
	}

	//relationTraits := body["relation-traits"].([]map[string]interface{})
	//if !ok {
	//	return
	//}

	tm = js.M{
		"value": "Application",
		"kids": js.S{
			js.M{
				"value": "Meta",
				"kids": js.S{
					js.M{
						"value": "Traits",
						"kids":  traitsToTreeModel(body),
					},
					js.M{
						"value": "Relations",
						"kids":  relationsToTreeModel(body),
					},
				},
			},
			js.M{
				"value": "Objects",
				"kids":  js.S{},
			},
		},
	}
	return
}

func fillToolBar(toolBar *ToolBar) {
	addToolBarImage(toolBar, "images/24/gtk-refresh.png", "tbRefresh")
	addToolBarImage(toolBar, "images/24/gnome-logout.png", "tbLogout")
}

func addToolBarImage(toolBar *ToolBar, img string, id string) {
	p := toolBar.AddImage(img)
	p.SetID(id)
	//p.Object().Set("fired", f)
}

func dispatchToolBarEvent(id string) {
	switch id {
	case "tbRefresh":
		go refreshMeta()
	}
}

func traitsToTreeModel(body map[string]interface{}) (model js.S) {
	//relations, ok := body["relations"].([]interface{})
	//if !ok {
	//	return
	//}
	traits, ok := body["traits"].([]interface{})
	if !ok {
		return
	}

	model = js.S{}
	for _, v := range traits {
		item := js.M{
			"value": v.(map[string]interface{})["name"],
		}
		model = append(model, item)
	}
	return
}

func relationsToTreeModel(body map[string]interface{}) (model js.S) {
	relations, ok := body["relations"].([]interface{})
	if !ok {
		return
	}
	//traits, ok := body["traits"].([]interface{})
	//if !ok {
	//	return
	//}
	model = js.S{}
	for _, v := range relations {
		item := js.M{
			"value": v.(map[string]interface{})["name"],
		}
		model = append(model, item)
	}
	return
}
