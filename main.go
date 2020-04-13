package main

import "log"

var zebkit *Zebkit
var innerWidth int
var innerHeight int

var session *Session

func init() {
	initResize()
	initSession("admin", "admin")

	zebkit = NewZebkit()
	zebkit.UIConfig("theme", "light")
}

func main() {
	err := testToken()
	if err != nil {
		err := login()
		if err != nil {
			log.Printf("Login error: %v", err)
		} else {
			log.Printf("Token: %v", session.Data.Token)
		}
	}
	zebkit.Require("ui", "layout", "data", BuildUI)
}
