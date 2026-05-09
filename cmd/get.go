package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)



var getCmd = &cobra.Command{
	Use:   "get [CMD]",
	Short: "Retreive information from the server",
	Long:  `Get information from the server.

	This should be used with the commands "user" and "messages".
	It should be run frequently if using mescli as CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
