package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Marks a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		var taskIds []int

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the arg: " + arg)
				continue
			}

			taskIds = append(taskIds, id)
		}

		fmt.Println(taskIds)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
