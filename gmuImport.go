// Ingest of MicroStrategy Exports
package schedules

import (
	"errors"
	"fmt"
	"github.com/ocelotsloth/xlsx"
	"os"
	"strings"
)

const sheetName = "Sheet1"

func GmuImport(files []string) {
	initLogging()
	startDatabase()
	defer database.Close()

	for _, file := range files {
		importFile(file)
	}
}

func locateHeader(xlFile *xlsx.File) (int, error) {
	for i := 1; i <= 40; i++ {
		currentCell := xlFile.Sheets[0].Cell(i, 1)
		currentValue := currentCell.String()
		if currentValue == "COURSE" {
			return i, nil
		}
	}
	return 0, errors.New("Valid header not found")
}

// Converts the different strings representing nil values from the import to
// blank strings. (A utility function)
func blankVarients(this string) string {
	if this == " - " || this == "...None" || this == "0 0" {
		return ""
	} else {
		return this
	}
}

func importFile(file string) {
	log.Infof("Opening excel file: %s", file)
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		log.Criticalf("File \"%s\" failed to open!", file)
		log.Critical(err)
		os.Exit(1)
	}
	log.Infof("%s opened successfully.", file)
	headerRow, err := locateHeader(xlFile)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	log.Infof("Header row is: %d", headerRow)
	// Begin loading in data
	sheet := xlFile.Sheet["Course Search - All Sections"]
	currentRow := headerRow + 1
	currentCourse := ""
	currentSection := ""
	tx := database.Begin()
	log.Info("Scanning rows...")
	for sheet.Cell(currentRow, 1).String() != "Total" {
		// Check if on the "Total" row (which should be skipped)
		if sheet.Cell(currentRow, 2).String() == "Total" {
			currentRow++
			continue
		}
		// Check if the current course has changed
		if sheet.Cell(currentRow, 1).String() != "" {
			currentCourse = sheet.Cell(currentRow, 1).String()
		}
		// Get data for the section
		sectionFullID := sheet.Cell(currentRow, 2).String()
		// Check for annoying inconsistencies
		if sectionFullID == "" {
			sectionFullID = currentSection
		} else {
			currentSection = sectionFullID
		}
		sectionID := strings.Split(sectionFullID, " ")[2]
		// Get the CRN
		crnInt, err := sheet.Cell(currentRow, 6).Int()
		if err != nil {
			log.Errorf("Row %d (crnInt): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the Room Capacity
		roomCap, err := sheet.Cell(currentRow, 27).Int()
		if err != nil {
			if sheet.Cell(currentRow, 27).String() == "" {
				roomCap = 999
			} else {
				log.Errorf("Row %d (roomCap): %s", currentRow, err)
				currentRow++
				continue
			}
		}
		// Get the X-List Property and convert to bool
		XlistStr := sheet.Cell(currentRow, 30).String()
		var Xlist bool
		if XlistStr == "Y" || XlistStr == "y" {
			Xlist = true
		} else if XlistStr == "N" || XlistStr == "n" {
			Xlist = false
		} else {
			log.Errorf("Row %d (Xlist): Invalid Value (%s)", currentRow, XlistStr)
			currentRow++
			continue
		}
		// Get the MinCreditHours
		MinCreditHoursInt, err := sheet.Cell(currentRow, 38).Int()
		if err != nil {
			log.Errorf("Row %d (MinCreditHours): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the Enrollment Limit
		EnrollmentLimitInt, err := sheet.Cell(currentRow, 39).Int()
		if err != nil {
			log.Errorf("Row %d (EnrollmentLimit): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the Enrolled Count
		EnrolledCount, err := sheet.Cell(currentRow, 40).Int()
		if err != nil {
			log.Errorf("Row %d (EnrolledCount): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the number of Seats Available
		SeatsAvail, err := sheet.Cell(currentRow, 41).Int()
		if err != nil {
			log.Errorf("Row %d (SeatsAvail): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the Waitlist Count
		WaitlistCount, err := sheet.Cell(currentRow, 42).Int()
		if err != nil {
			log.Errorf("Row %d (WaitlistCount): %s", currentRow, err)
			currentRow++
			continue
		}
		// Get the data for the Sessions
		days := blankVarients(sheet.Cell(currentRow, 22).String())
		times := blankVarients(sheet.Cell(currentRow, 23).String())
		secDays := blankVarients(sheet.Cell(currentRow, 33).String())
		secTimes := blankVarients(sheet.Cell(currentRow, 34).String())
		tertDays := blankVarients(sheet.Cell(currentRow, 35).String())
		tertTimes := blankVarients(sheet.Cell(currentRow, 36).String())
		var SessionsData = []struct {
			days, times string
		}{
			{days, times},
			{secDays, secTimes},
			{tertDays, tertTimes},
		}
		var Sessions []Session
		for _, SessionData := range SessionsData {
			var nextSession Session
			// Check to see if we even need to process this iteration
			if SessionData.days == "" || SessionData.times == "" {
				continue
			}
			for _, character := range SessionData.days {
				switch character {
				case 'M':
					nextSession.Mon = true
				case 'T':
					nextSession.Tues = true
				case 'W':
					nextSession.Wed = true
				case 'R':
					nextSession.Thurs = true
				case 'F':
					nextSession.Fri = true
				case 'S':
					nextSession.Sat = true
				case 'U':
					nextSession.Sun = true
				default:
					log.Error(fmt.Sprintf("Invalid Character (%c) found.", character))
				}
			}
			// Parse Time
			nextSession.Start, nextSession.End, err = parseTimeRange(SessionData.times)
			if err != nil {
				log.Error(err)
				continue
			}
			// Add Location
			nextSession.Location = blankVarients(sheet.Cell(currentRow, 25).String())
			Sessions = append(Sessions, nextSession)
		}

		// Construct the newSection
		newSection := Section{
			Semester:     file,
			Course:       currentCourse,
			Section:      sectionID,
			CRN:          crnInt,
			LinkCode:     sheet.Cell(currentRow, 7).String(),
			ScheduleType: sheet.Cell(currentRow, 8).String(),
			SectionTitle: sheet.Cell(currentRow, 11).String(),
			Instructor:   sheet.Cell(currentRow, 16).String(),
			StartDate:    sheet.Cell(currentRow, 18).String(),
			EndDate:      sheet.Cell(currentRow, 21).String(),
			Sessions:     Sessions,
			Days:         days,
			Times:        times,
			SecDays:      secDays,
			SecTimes:     secTimes,
			TertDays:     tertDays,
			TertTimes:    tertTimes,
			Location:     blankVarients(sheet.Cell(currentRow, 25).String()),
			RoomCapacity: roomCap,
			Status:       sheet.Cell(currentRow, 28).String(),
			CampusName:   sheet.Cell(currentRow, 29).String(),
			XList:        Xlist,
			XListNumber:  sheet.Cell(currentRow, 31).String(),
			SectionNotes: sheet.Cell(currentRow, 32).String(),
			MinCRHours:   MinCreditHoursInt,
			Limit:        EnrollmentLimitInt,
			Enrolled:     EnrolledCount,
			SeatsAvail:   SeatsAvail,
			Waitlist:     WaitlistCount,
		}
		tx.Create(&newSection)
		currentRow++
	}
	log.Info("Commiting database changes...")
	tx.Commit()
	log.Info("Done importing.")
}
