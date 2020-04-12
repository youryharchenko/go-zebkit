package main

var zebkit *Zebkit
var innerWidth int
var innerHeight int

func init() {
	initResize()
	zebkit = NewZebkit()
}

func main() {
	zebkit.UIConfig("theme", "light")
	zebkit.Require("ui", "layout", BuildUI)
}
