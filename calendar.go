package schedules

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Calendar Generation

func HandleCalendarGET(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Calendar Generate")
}
