package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "mescli",
	Short: "An end-to-end messaging service for the terminal",
	Long:  `mescli is an end-to-end messaging service built for the terminal

	As the name indicates, mescli is intended as a CLI to send quick 
	encrypted messages. As the development became more complex, it 
	was decided to include as well a TUI interface.

	The primary feature of mescli is the asynchronous encryption that is 
	achieved using the Signal encryption protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will be the point when the TUI will open.")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
