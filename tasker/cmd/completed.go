package cmd

import (
	"fmt"
	"tasker/db"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Lists all completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.ConnectDB()
		if err != nil {
			panic(err)
		}

		store := db.NewStore(database)

		tasks, err := store.FindAllCompletedTasks()
		if err != nil {
			panic(err)
		}

		if len(tasks) == 0 {
			fmt.Println("You haven't completed any tasks yet :(")
			return
		}

		fmt.Println("Here are your completed tasks:")
		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
