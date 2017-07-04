package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Course Collections

func HandleCourseIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var courses []Course
	err := database.Find(&courses).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Courses\n")
		return
	}
	json.NewEncoder(rw).Encode(courses)
}

func HandleCourseCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusCreate")
}

// Course Singular

func HandleCourseShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleCourseEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleCourseDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
