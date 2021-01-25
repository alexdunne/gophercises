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

		database, err := db.ConnectDB()
		if err != nil {
			panic(err)
		}

		store := db.NewStore(database)

		_, err = store.CreateTask(&db.Task{Description: task})
		if err != nil {
			panic(err)
		}

		fmt.Println("Task added")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
