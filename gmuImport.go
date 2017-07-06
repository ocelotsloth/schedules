// Ingest of MicroStrategy Exports
package schedules

import (
	"errors"
	"github.com/tealeg/xlsx"
	"os"
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

func importFile(file string) {
	log.Infof("Opening excel file: %s", file)
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		log.Criticalf("File \"%s\" failed to open!", file)
		log.Critical(err)
		os.Exit(1)
	}

	// currently hangs right here....not sure why

	log.Infof("%s opened successfully.", file)

	headerRow, err := locateHeader(xlFile)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	log.Infof("Header row is: %d", headerRow)
}

func locateHeader(xlFile *xlsx.File) (int, error) {
	for i := 1; i <= 40; i++ {
		currentCell := xlFile.Sheets[0].Cell(i, 1)
		currentValue, err := currentCell.String()
		if err != nil {
			return 0, err
		}
		log.Infof("| %d | %s |", i, currentValue)
		if currentValue == "COURSE" {
			return i, nil
		}
	}
	return 0, errors.New("Valid header not found")
}
