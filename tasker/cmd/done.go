package cmd

import (
	"fmt"
	"strconv"
	"tasker/db"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Marks a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		var taskIDs []int

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the arg: " + arg)
				continue
			}

			taskIDs = append(taskIDs, id)
		}

		database, err := db.ConnectDB()
		if err != nil {
			panic(err)
		}

		store := db.NewStore(database)

		for _, taskID := range taskIDs {
			store.DeleteTask(taskID)
			fmt.Printf("Task %d marked as complete\n", taskID)
		}
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
