package cmd

import (
	"fmt"
	"strconv"
	"tasker/db"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a task",
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
			fmt.Printf("Task %d removed\n", taskID)
		}
	},
}

var rmAlias = &cobra.Command{
	Use:   "rm",
	Short: "alias for remove",
	Run: func(cmd *cobra.Command, args []string) {
		removeCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(rmAlias)
}
