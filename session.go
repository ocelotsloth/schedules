package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Session Collections

func HandleSessionIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var sessions []Session
	err := database.Find(&sessions).Error
	if err != nil {
		log.Critical(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error Retrieving Sessions\n")
		return
	}
	json.NewEncoder(rw).Encode(sessions)
}

func HandleSessionCreate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "CampusCreate")
}

// Session Singular

func HandleSessionShow(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleSessionEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func HandleSessionDelete(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
