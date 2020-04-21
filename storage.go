package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
)

// Local -
type Local struct {
	Obj  *js.Object
	Data LocalData
}

// LocalData -
type LocalData struct {
	User       string    `json:"user"`
	Secret     string    `json:"secret"`
	ObjQueries []*Query  `json:"obj-queries"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

// Query -
type Query struct {
	Name   string                 `json:"name"`
	Path   string                 `json:"path"`
	View   map[string]interface{} `json:"view"`
	Source map[string]interface{} `json:"source"`
}

// NewLocal -
func NewLocal(data LocalData) (s *Local) {
	data.Created = time.Now()
	data.Updated = time.Now()

	s = &Local{
		Obj:  js.Global.Get("localStorage"),
		Data: data,
	}
	s.Save()
	return
}

// OpenLocal -
func OpenLocal() (s *Local) {
	o := js.Global.Get("localStorage")
	d := o.Get("localData").String()
	log.Printf("Get localData:%s.", d)
	sd := LocalData{}
	if len(d) > 0 && d != "undefined" {
		err := json.Unmarshal([]byte(d), &sd)
		if err != nil {
			log.Println(err)
			return
		}
		s = &Local{
			Obj:  o,
			Data: sd,
		}
		s.Save()
	} else {
		sd.Created = time.Now()
		s = &Local{
			Obj:  o,
			Data: sd,
		}
		s.Save()
	}
	return
}

// Save -
func (s *Local) Save() {
	s.Data.Updated = time.Now()
	d, err := json.Marshal(s.Data)
	if err != nil {
		log.Print(err.Error())
	}
	s.Obj.Set("localData", string(d))
}

// Session -
type Session struct {
	Obj  *js.Object
	Data SessionData
}

// SessionData -
type SessionData struct {
	User    string    `json:"user"`
	Secret  string    `json:"secret"`
	UUID    string    `json:"uuid"`
	Token   string    `json:"token"`
	Meta    js.M      `json:"meta"`
	Item    string    `json:"item"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// NewSession -
func NewSession(data SessionData) (s *Session) {
	data.Created = time.Now()
	data.Updated = time.Now()
	data.UUID = uuid.New().String()
	s = &Session{
		Obj:  js.Global.Get("sessionStorage"),
		Data: data,
	}
	s.Save()
	return
}

// OpenSession -
func OpenSession() (s *Session) {
	o := js.Global.Get("sessionStorage")
	d := o.Get("sessionData").String()
	log.Printf("Get sessionData:%s.", d)
	if len(d) > 0 && d != "undefined" {
		sd := SessionData{}
		err := json.Unmarshal([]byte(d), &sd)
		if err != nil {
			log.Println(err)
			return
		}
		s = &Session{
			Obj:  o,
			Data: sd,
		}
		s.Save()
	}
	return
}

// Save -
func (s *Session) Save() {
	s.Data.Updated = time.Now()
	d, err := json.Marshal(s.Data)
	if err != nil {
		log.Print(err.Error())
	}
	s.Obj.Set("sessionData", string(d))
}
