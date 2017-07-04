// Ingest of MicroStrategy Exports
package schedules

import (
	"fmt"
	"github.com/ocelotsloth/phantomjs"
	"os"
	"net/http"
)

func GmuFetch() {
	// Start the process once.
	phantomjs.DefaultProcess.CmdArgs = []string{"--ignore-ssl-errors=true"}
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer phantomjs.DefaultProcess.Close()

	// Create a web page
	page, err := phantomjs.CreateWebPage()
	if err != nil {
		panic(err)
	}
	defer page.Close()

	loginCookie := http.Cookie{
		Name: "JSESSIONID",
		Value: "B6B993634001C6BBB27A13970A8C9381",
		Path: "/MicroStrategy/",
		Domain: "microstrategy.gmu.edu",
	}
	page.SetCookies([]*http.Cookie{&loginCookie})

	// Open a URL.
	if err := page.Open("https://microstrategy.gmu.edu/MicroStrategy/servlet/mstrWeb?evt=2048001&evt=1024001&src=mstrWeb.2048001&src=mstrWeb.shared.fbb.1024001&documentID=B538AE6F43A7C5E2A517B58C478B3497&currentViewMedia=1&visMode=0&events=-5008*.mstrWeb***.shared***.5008*.name*.ltH*.value*.589.5008*.mstrWeb***.shared***.5008*.name*.ltW*.value*.200.5008*.mstrWeb***.shared***.5008*.name*.ltH*.value*.589.5008*.mstrWeb***.shared***.5008*.name*.ltW*.value*.200.5008*.mstrWeb***.shared***.5008*.name*.ltH*.value*.589.5008*.mstrWeb***.shared***.5008*.name*.ltW*.value*.200.5008*.mstrWeb***.shared***.5008*.name*.ltH*.value*.589.5008*.mstrWeb***.shared***.5008*.name*.ltW*.value*.200_&1024001=1&evtorder=1024001&2048001=1&mstrWeb=-*-lfj*_rEh4TsVKYAj9SHqSqOP9bw4%3D.GMU+Public+Reports.*-CJ8PNnqytGuuV2mN_&shared=*-1.*-1.0.0.0&smartBanner=*0.*-1.0&ftb=0.8DE367974292BA099CA0E2A1562DB9EB.*0.8.0.0-8.18_268453447.*-1.1.*0&fb=0.8DE367974292BA099CA0E2A1562DB9EB.Course%2BEnrollment.8.0.0-8.768.769.774.770.773.775.55.776.777.779.72_268453447.*-1.1.*0&Server=INDICIUM.GMU.EDU&Port=34952&Project=GMU+Public+Reports&"); err != nil {
		fmt.Println(err)
		panic(err)
	}

	// setup viewport and render
	if err := page.SetViewportSize(1024, 768); err != nil {
		panic(err)
	}
	if err := page.Render("test.png", "png", 100); err != nil {
		panic(err)
	}
}
