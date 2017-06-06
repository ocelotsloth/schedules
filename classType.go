package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Class Type Collections

func HandleClassTypeIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var classTypes []ClassType
	err := database.Find(&classTypes).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Class Types\n")
		return
	}
	json.NewEncoder(rw).Encode(classTypes)
}

func HandleClassTypeCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusCreate")
}

// Class Type Singular

func HandleClassTypeShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleClassTypeEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleClassTypeDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
