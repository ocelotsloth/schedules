package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var log = logging.MustGetLogger("schedules_server")
var database *gorm.DB

func main() {
	initLogging()
	startDatabase()
	defer database.Close()
	startAPI()
}

// initLogging Configures and initializes logging for the daemon
func initLogging() {

	// Example format string. Everything except the message has a custom color
	// which is dependent on the log level. Many fields have a custom output
	// formatting too, eg. the time returns the hour down to the milli second.
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

// startDatabase Starts connection to database
func startDatabase() {
	log.Info("Connecting to database...")
	db, err := gorm.Open("sqlite3", "./schedules.db")
	database = db
	database.AutoMigrate(Campus{})
	database.AutoMigrate(Semester{})
	database.AutoMigrate(Course{})
	database.AutoMigrate(Session{})
	database.AutoMigrate(ClassType{})
	database.AutoMigrate(Section{})
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %s\n", err.Error())
		os.Exit(1)
	}
	log.Info("Connected to Database!")

	// Check that the class types are in the database
	var classTypes []ClassType
	db.Find(&classTypes)
	if len(classTypes) == 0 {
		log.Info("Adding Initial Class Types...")

		db.Save(&ClassType{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Lecture",
		})
		db.Save(&ClassType{
			Model: gorm.Model{
				ID: 2,
			},
			Name: "Lab",
		})
		db.Save(&ClassType{
			Model: gorm.Model{
				ID: 3,
			},
			Name: "Recitation",
		})
		db.Save(&ClassType{
			Model: gorm.Model{
				ID: 4,
			},
			Name: "Other",
		})
		log.Info("Class Types Added!")
	}

}

var templates *template.Template

func startAPI() {
	loadTemplates()
	r := httprouter.New()

	// Home
	r.GET("/", HandleHomeGET)

	// Campus Collection
	r.GET("/campuses", HandleCampusIndex)
	r.POST("/campuses", HandleCampusCreate)

	// Campus Singular
	r.GET("/campuses/:slug", HandleCampusShow)
	r.GET("/campuses/:slug/edit", HandleCampusEdit)
	r.DELETE("/campuses/:slug", HandleCampusDelete)

	// Semester Collection
	r.GET("/semesters", HandleSemesterIndex)
	r.POST("/semesters", HandleSemesterCreate)

	// Semester Singular
	r.GET("/semesters/:slug", HandleSemesterShow)
	r.GET("/semesters/:slug/edit", HandleSemesterEdit)
	r.DELETE("/semesters/:slug", HandleSemesterDelete)

	// Course Collection
	r.GET("/courses", HandleCourseIndex)
	r.POST("/courses", HandleCourseCreate)
	
	// Course Singular
	r.GET("/courses/:slug", HandleCourseShow)
	r.GET("/courses/:slug/edit", HandleCourseEdit)
	r.DELETE("/courses/:slug", HandleCourseDelete)

	// Session Collection
	r.GET("/sessions", HandleSessionIndex)
	r.POST("/sessions", HandleSessionCreate)

	// Session Singular
	r.GET("/sessions/:slug", HandleSessionShow)
	r.GET("/sessions/:slug/edit", HandleSessionEdit)
	r.DELETE("/sessions/:slug", HandleSessionDelete)

	// ClassType Collection
	r.GET("/classTypes", HandleClassTypeIndex)
	r.POST("/classTypes", HandleClassTypeCreate)

	// ClassType Singular
	r.GET("/classTypes/:slug", HandleClassTypeShow)
	r.GET("/classTypes/:slug/edit", HandleClassTypeEdit)
	r.DELETE("/ClassTypes/:slug", HandleClassTypeDelete)

	// Section Collection
	r.GET("/sections", HandleSectionIndex)
	r.PUT("/sections", HandleSectionCreate)

	// Section Singular
	r.GET("/sections/:slug", HandleSectionShow)
	r.GET("/sections/:slug/edit", HandleSectionEdit)
	r.DELETE("/sections/:slug", HandleSectionDelete)
	
	// Calendar Generation
	r.GET("/calendar", HandleCalendarGET)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)

	http.ListenAndServe(":8080", n)
}

func loadTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			log.Infof("Loaded Template: %s\n", filename)
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}
	templates, err = template.New("").ParseFiles(allFiles...) //parses all .html files in the 'templates' folder

}
