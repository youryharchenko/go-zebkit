package main

var zebkit *Zebkit
var innerWidth int
var innerHeight int

var session *Session
var local *Local

func init() {
	initResize()
	initLocal()
	initSession()

	zebkit = NewZebkit()
	zebkit.UIConfig("theme", "light")
}

func main() {
	zebkit.Require("ui", "layout", "data", "draw", BuildUI)
}
