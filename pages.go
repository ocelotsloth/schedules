package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleHomeGET(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplate(rw, "home", nil)
}
