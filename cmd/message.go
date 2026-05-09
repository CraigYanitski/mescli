package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)



var getMessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Get messages from the server",
	Long:  `Get messages stored on the server.

	This requires the user to be logged in.
	Any messages sent to the user are stored on the database until 
	retrieved. The messages will be displayed by user then by time.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

var viewMessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "View messages from local disk",
	Long:  `View messages from local disk.

	The user email or UUID may be specified for a specific conversation.
	This prints the five most recent messages in from the 3 most recent conversations.
	It will soon be possible to alter these numbers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This command is not yet implemented...")
	},
}

func init() {
	getCmd.AddCommand(getMessagesCmd)
	viewCmd.AddCommand(viewMessagesCmd)
}

