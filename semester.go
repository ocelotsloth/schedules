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
	// https://stackoverflow.com/a/15685432
	var newSemester Semester
	err := json.NewDecoder(r.Body).Decode(&newSemester)
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error Decoding Request\n")
		return
	}
	defer r.Body.Close()

	tx := database.Begin()
	if err := tx.Create(&newSemester).Error; err != nil {
		tx.Rollback()
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error writing to database\n")
		return
	}
	tx.Commit()

	log.Info("Created new Semester with ID=" + string(newSemester.ID) + ".")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, newSemester.ID)
}

// Semester Singular

func HandleSemesterShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var semester Semester
	dbOp := database.Where("Slug = ?", p.ByName("Slug")).First(&semester)
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
	json.NewEncoder(rw).Encode(semester)
}

func HandleSemesterEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleSemesterDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Search for Campus Entry
	var semester Semester
	dbOp := database.Where("Slug = ?", p.ByName("Slug")).First(&semester)
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
	err := tx.Unscoped().Delete(&semester).Error
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Internal Server Error\n")
		log.Critical(dbOp.Error.Error())
		return
	}
	tx.Commit()
	// Send a response back
	log.Info(fmt.Sprintf("Record \"%s\" deleted.", semester.Slug))
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Record Successfully Deleted.")
}
