package cmd

import (
	"fmt"
	"github.com/ocelotsloth/schedules"
	"github.com/spf13/cobra"
	"os"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "import [file.xlsx]...",
	Short: "Imports the given excel data files to the database.",
	Long: `This command loads the data located in a microstrategy.gmu.edu
report and imports it into the database for the server to use
immediately.

This command can be run while an active server instance is
also running, and will have the effect of instantly updating
the data.

To ensure that partially loaded data is not sent back to end
users, the entire database operation is staged before being
applied in one fell swoop.

---

To grab a datafile, visit microstrategy.gmu.edu, login as a
guest user, and navigate to the course enrollment report.
You will want to export this report for a specific semester.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			if _, err := os.Stat(arg); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("'%s' is not a valid file.\n", arg))
				os.Exit(-43)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd Run with args: %v\n", args)
		schedules.GmuImport(args)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags ahttps://github.com/magit/magit/wiki/Emacscliend confighttps://github.com/magit/magit/wiki/Emacsclieuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
