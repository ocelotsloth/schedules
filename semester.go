package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Semester Collections

func HandleSemesterIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var semesters []Semester
	err := database.Find(&semesters).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Semesters\n")
		return
	}
	json.NewEncoder(rw).Encode(semesters)
}

func HandleSemesterCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusCreate")
}

// Semester Singular

func HandleSemesterShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleSemesterEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleSemesterDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
