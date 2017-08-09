package schedules

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Section Collections

// Queries the database for a particular section in a given semester.
func GetSectionByCRN(crn int, semester string) (Section, error) {
	if database == nil {
		return Section{}, errors.New("Database not yet initialized")
	}
	var response []Section
	database.Where(&Section{CRN: crn, Semester: semester}).Find(&response)
	if len(response) > 1 {
		return Section{}, errors.New("Multiple matches found")
	}
	switch length := len(response); {
	case length > 1:
		return Section{}, errors.New(fmt.Sprintf("Multiple matches found with crn: %d", crn))
	case length == 0:
		return Section{}, errors.New(fmt.Sprintf("Section %d does not exist", crn))
	}
	return response[0], nil
}

// Queries the database for a list of sections in a given semester.
func GetSectionsByCRN(crns []int, semester string) ([]Section, error) {
	var sections []Section
	for _, crn := range crns {
		newSection, err := GetSectionByCRN(crn, semester)
		if err != nil {
			return []Section{}, err
		}
		sections = append(sections, newSection)
	}
	return sections, nil
}

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
