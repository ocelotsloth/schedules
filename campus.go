package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

// Campus Collections

func HandleCampusIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var campuses []Campus
	err := database.Find(&campuses).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Campuses\n")
		return
	}
	json.NewEncoder(rw).Encode(campuses)
}

func HandleCampusCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// https://stackoverflow.com/a/15685432
	decoder := json.NewDecoder(r.Body)
	var newCampus Campus
	err := decoder.Decode(&newCampus)
	if err != nil {
		log.Critical(err.Error())
		panic(err)
	}
	defer r.Body.Close()

	if database.NewRecord(newCampus) {
		database.Create(&newCampus)
	}
	log.Info("Created new campus with ID=" + string(newCampus.ID) + ".")

	fmt.Fprintln(rw, newCampus.ID)
}

// Campus Singular

func HandleCampusShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusShow")
}

func HandleCampusEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusEdit")
}

func HandleCampusDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusDelete")
}

