package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)



var viewCmd = &cobra.Command{
	Use:   "view [CMD]",
	Short: "View information from local disk",
	Long:  `View information from local disk.

	This should be used with the commands "users" and "messages".
	It will be a convenient method display local information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
