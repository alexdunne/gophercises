package cmd

import (
	"fmt"
	"tasker/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all active tasks",
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.ConnectDB()
		if err != nil {
			panic(err)
		}

		store := db.NewStore(database)

		tasks, err := store.FindAllTasks()
		if err != nil {
			panic(err)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks :D")
			return
		}

		fmt.Println("You have the following tasks:")
		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.Key, task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
