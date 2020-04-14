package main

var zebkit *Zebkit
var innerWidth int
var innerHeight int

var session *Session

func init() {
	initResize()
	initSession("admin", "")

	zebkit = NewZebkit()
	zebkit.UIConfig("theme", "light")
}

func main() {
	zebkit.Require("ui", "layout", "data", BuildUI)
}
