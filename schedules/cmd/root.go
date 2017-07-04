package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "schedules",
		Short: "A web service for generating class iCal files",
		Long:  `Schedules is a server application that allows students to generate iCal files automatically which contain their entire class schedule.`,
	}
)

func Execute() {
	rootCmd.Execute()
}

// Init function runs automatically before Excecute()
func init() {

}
