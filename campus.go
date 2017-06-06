package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(campuses)
}

func HandleCampusCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// https://stackoverflow.com/a/15685432
	var newCampus Campus
	err := json.NewDecoder(r.Body).Decode(&newCampus)
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error Decoding Request\n")
		return
	}
	defer r.Body.Close()

	tx := database.Begin()
	if err := tx.Create(&newCampus).Error; err != nil {
		tx.Rollback()
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error writing to database\n")
		return
	}
	tx.Commit()

	log.Info("Created new campus with ID=" + string(newCampus.ID) + ".")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, newCampus.ID)
}

// Campus Singular

func HandleCampusShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var campus Campus
	dbOp := database.Where("Slug = ?", p.ByName("Slug")).First(&campus)
	if dbOp.RecordNotFound() {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Record not found.\n")
		return
	}
	if dbOp.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Internal Server Error\n")
		log.Critical(dbOp.Error.Error())
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(campus)
}

// TODO: MAKE WORK
func HandleCampusEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// https://stackoverflow.com/a/15685432
	var campus Campus
	var changedCampus interface{}
	err := json.NewDecoder(r.Body).Decode(&changedCampus)
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error Decoding Request\n")
		return
	}
	defer r.Body.Close()

	dbOp := database.Where("Slug = ?", p.ByName("Slug")).First(&campus)
	if dbOp.RecordNotFound() {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Record not found.\n")
		return
	}
	if dbOp.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Internal Server Error\n")
		log.Critical(dbOp.Error.Error())
		return
	}

	// mergo.Merge(&changedCampus, campus)

	json.NewEncoder(rw).Encode(changedCampus)
}

func HandleCampusDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Search for Campus Entry
	var campus Campus
	dbOp := database.Where("Slug = ?", p.ByName("Slug")).First(&campus)
	if dbOp.RecordNotFound() {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "Record not found.\n")
		return
	}
	if dbOp.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Internal Server Error\n")
		log.Critical(dbOp.Error.Error())
		return
	}
	// If found, delete it.
	tx := database.Begin()
	err := tx.Unscoped().Delete(&campus).Error
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Internal Server Error\n")
		log.Critical(dbOp.Error.Error())
		return
	}
	tx.Commit()
	// Send a response back
	log.Info(fmt.Sprintf("Record \"%s\" deleted.", campus.Slug))
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Record Successfully Deleted.")
}
