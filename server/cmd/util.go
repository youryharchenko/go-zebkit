package main

import (
	"bytes"
	"html/template"
	"io/ioutil"

	acc "github.com/youryharchenko/accountancy"
)

func initAll(db string, query string) (response string, err error) {
	tmplFile, err := ioutil.ReadFile(query + "/init.json")
	if err != nil {
		return
	}
	tmpl, err := template.New("data").Parse(string(tmplFile))
	if err != nil {
		return
	}
	data := map[string]interface{}{
		"db": db,
	}
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, data)
	if err != nil {
		return
	}
	response, err = acc.RunBatch(buff.String(), nil, nil)
	if err != nil {
		return
	}
	return
}
