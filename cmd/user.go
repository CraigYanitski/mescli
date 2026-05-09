package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)



var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user information from the server",
	Long:  `Get user information from the server.

	The user email must be specified.
	This prints the user's UUID.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

var viewUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "View user information from local disk",
	Long:  `View user information from local disk.

	The user email may be specified for specific information.
	This prints each user's UUID, message statistics, and most recent message.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

func init() {
	getCmd.AddCommand(getUserCmd)
	viewCmd.AddCommand(viewUsersCmd)
}

