package main

var zebkit *Zebkit

func init() {

	zebkit = NewZebkit()
	//log.Println("zebkit", zebkit)

}

func main() {
	zebkit.UIConfig("theme", "light")
	zebkit.Require("ui", "layout", BuildUI)
}
