package main

import "log"

func menuList() (list []interface{}) {
	list = []interface{}{
		map[string]interface{}{
			"content": "File",
			"sub": []interface{}{
				map[string]interface{}{
					"content": "New",
					"id":      "menuNew",
					"handler": menuHandlerNew,
				},
				map[string]interface{}{
					"content": "Open",
					"id":      "menuOpen",
					"handler": menuHandlerOpen,
				},
				"-", // line
				map[string]interface{}{
					"content": "Exit",
					"id":      "menuExit",
					"handler": menuHandlerExit,
				},
			},
		},
		map[string]interface{}{
			"content": "Edit",
			"sub": []interface{}{
				map[string]interface{}{
					"content": "Cut",
					"id":      "menuCut",
					"handler": menuHandler,
				},
				map[string]interface{}{
					"content": "Copy",
					"id":      "menuCopy",
					"handler": menuHandler,
				},
				map[string]interface{}{
					"content": "Paste",
					"id":      "menuPaste",
					"handler": menuHandler,
				},
			},
		},
		map[string]interface{}{
			"content": "Help",
			"sub": []interface{}{
				map[string]interface{}{
					"content": "About",
					"id":      "menuAbout",
					"handler": menuHandler,
				},
			},
		},
	}
	return
}

func menuHandler() {
	log.Println("Menu Click!")
}

func menuHandlerNew() {
	log.Println("New Click!")
}

func menuHandlerOpen() {
	log.Println("Open Click!")
}

func menuHandlerExit() {
	log.Println("Exit Click!")
}
