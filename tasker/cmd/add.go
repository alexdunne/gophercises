package cmd

import (
	"fmt"
	"strings"
	"tasker/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		fmt.Println(task)

		database, err := db.ConnectDB()
		if err != nil {
			panic(err)
		}

		h := NewCommandHandler(database)

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
