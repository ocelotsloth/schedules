package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

// Section Collections

func HandleSectionIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var sections []Section
	err := database.Find(&sections).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Sections\n")
		return
	}
	json.NewEncoder(rw).Encode(sections)
}

func HandleSectionCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Section Create")
}

// Section Singular

func HandleSectionShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Section Show")	
}

func HandleSectionEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Section Edit")
}

func HandleSectionDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Section Delete")
}

