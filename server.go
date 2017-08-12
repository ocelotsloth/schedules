package schedules

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
	"os/user"
	"path/filepath"
	"strings"
)

var log = logging.MustGetLogger("schedules_server")
var database *gorm.DB

func StartServer() {
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

	// Get Database Path
	usr, _ := user.Current()
	dir := usr.HomeDir
	// Change later so that this comes from Viper
	dbPath := filepath.Join(dir, "Desktop/schedules.db")
	log.Info(fmt.Sprintf("Database Path: %s", dbPath))

	// Open Database
	db, err := gorm.Open("sqlite3", dbPath)
	database = db
	db.AutoMigrate(Campus{})
	//db.AutoMigrate(Semester{})
	db.AutoMigrate(Course{})
	db.AutoMigrate(Session{})
	db.AutoMigrate(Section{})
	db.Exec("CREATE VIEW semesters AS SELECT DISTINCT semester AS 'name' FROM sections;")
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %s\n", err.Error())
		os.Exit(1)
	}
	log.Info("Connected to Database!")
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
	r.GET("/campuses/:Slug", HandleCampusShow)
	r.PATCH("/campuses/:Slug", HandleCampusEdit)
	r.DELETE("/campuses/:Slug", HandleCampusDelete)

	// Semester Collection
	r.GET("/semesters", HandleSemesterIndex)
	r.POST("/semesters", HandleSemesterCreate)

	// Semester Singular
	r.GET("/semesters/:Slug", HandleSemesterShow)
	r.PATCH("/semesters/:Slug", HandleSemesterEdit)
	r.DELETE("/semesters/:Slug", HandleSemesterDelete)

	// Course Collection
	r.GET("/courses", HandleCourseIndex)
	r.POST("/courses", HandleCourseCreate)

	// Course Singular
	r.GET("/courses/:Slug", HandleCourseShow)
	r.PATCH("/courses/:Slug", HandleCourseEdit)
	r.DELETE("/courses/:Slug", HandleCourseDelete)

	// Session Collection
	r.GET("/sessions", HandleSessionIndex)
	r.POST("/sessions", HandleSessionCreate)

	// Session Singular
	r.GET("/sessions/:Slug", HandleSessionShow)
	r.PATCH("/sessions/:Slug", HandleSessionEdit)
	r.DELETE("/sessions/:Slug", HandleSessionDelete)

	// Section Collection
	r.GET("/sections", HandleSectionIndex)
	r.PUT("/sections", HandleSectionCreate)

	// Section Singular
	r.GET("/sections/:Slug", HandleSectionShow)
	r.PATCH("/sections/:Slug", HandleSectionEdit)
	r.DELETE("/sections/:Slug", HandleSectionDelete)

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
