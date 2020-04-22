package main

import "github.com/youryharchenko/go-zebkit/zkit"

var zebkit *zkit.Zebkit
var innerWidth int
var innerHeight int

var session *Session
var local *Local

var rootCanvas *zkit.Panel

var formLogin *Form

var toolBar *zkit.ToolBar
var statusBar *zkit.StatusBarPan
var tree *zkit.Tree
var treeModel *zkit.TreeModel
var tabs *zkit.Tabs

func init() {
	initResize()
	initLocal()
	initSession()

	zebkit = zkit.NewZebkit()
	zebkit.UIConfig("theme", "light")
}

func main() {
	zebkit.Require("ui", "layout", "data", "draw", BuildUI)
}
