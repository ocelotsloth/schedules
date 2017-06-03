package main

import (
	//"github.com/julienschmidt/httprouter"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t := templates.Lookup(tmpl)
	//t, _ := template.ParseFiles("templates/"+tmpl + ".html")
	//t := template.New("fieldname example")
	//t, _ = t.Parse("hello {{.Title}}!")
	t.Execute(w, data)
}
