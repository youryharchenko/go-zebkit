package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	jsb "github.com/gopherjs/jsbuiltin"
	xhr "github.com/rocketlaunchr/gopherjs-xhr"
	"go4.org/sort"
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
	local.Data.ObjQueries = makeObjQueries()
	log.Printf("Open local storage")
}

func makeObjQueries() (qrs []*Query) {
	qrs = []*Query{
		{
			Name: "Test",
			Src: `{
				"request": {
					"command": "select",
					"service": "object"
				},
				"db": {
					"driver": "",
					"connection": "",
					"show": true
				},
				"body": {
					"filter": {
						"condition": "status = ?",
						"params": [0]
					},
					"orderBy": "name asc",
					"skip": 0,
					"limit": 50,
					"filterProps": "props.pay && props.pay.request.service == 'nparts'",
					"fields": "{'name':obj.name,'terminal':obj.props.pay.request.terminalId,'amount':obj.props.pay.request.body.amount,'flag':true}"
				}
			}`,
		},
	}
	return
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

func runQuery(src string) (response string, err error) {

	chOk := make(chan bool, 0)
	go testAndLogin(&rootCanvas.Layoutable, chOk)
	ok := <-chOk
	if !ok {
		err = errors.New("test and login failed")
		return
	}

	response, err = post("/api/run", src)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func sortItems(a []interface{}) (s []interface{}) {
	sort.Slice(a, func(i, j int) (r bool) {
		name1, ok := a[i].(map[string]interface{})["name"].(string)
		if !ok {
			name1 = "Noname"
		}
		name2, ok := a[j].(map[string]interface{})["name"].(string)
		if !ok {
			name1 = "Noname"
		}
		r = name1 < name2
		return // family[i].Age < family[j].Age
	})
	s = a
	return
}

/*
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
*/
func fillToolBar(toolBar *ToolBar) {
	addToolBarImage(toolBar, "images/24/gtk-refresh.png", "tbRefresh")
	addToolBarImage(toolBar, "images/24/stock_media-play.png", "tbRun")
	addToolBarImage(toolBar, "images/24/gnome-logout.png", "tbLogout")
}

func addToolBarImage(toolBar *ToolBar, img string, id string) {
	p := toolBar.AddImage(img)
	p.SetID(id)
	//p.Object().Set("fired", f)
}

/*
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
*/
func findItemInTreeModel(tm *TreeModel, find string) (item *Item) {
	var name, t string
	af := strings.Split(find, " ")
	if len(af) > 1 {
		name, t = af[0], af[1]
	}

	treeModel.Iterate(treeModel.Root(), func(i *js.Object) {
		v := NewItem(i).Value()
		if name == v.Get("name").String() && t == v.Get("type").String() {
			item = NewItem(i)
		}
	})

	if item != nil {
		log.Println("Found item:", item.Value().Get("name").String(), item.Value().Get("type").String())
	}

	return
}

func testAndLogin(root *Layoutable, ch chan bool) {
	err := testToken()
	if err != nil {
		i := 0
		for {
			i++
			log.Printf("Login retry %v", i)
			if i > 3 {
				log.Printf("Used max retries: %v", i)
				ch <- false
				return
			}
			log.Printf("Login - Will RunModal")
			res, err := formLogin.RunModal(root)
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

	ch <- true
}
