package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
)

// Sess

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
