package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save credentials in config file",
	Long:  `Login to the mescli service.

	This can be completed in either the CLI or TUI.
	It will need to be repeated if you logout or your JWT expires.`,
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("This command still needs to be implemented...")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove credentials from config file",
	Long:  `Logout of the mescli service.

	This can be completed in either the CLI or TUI.`,
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("This command still needs to be implemented...")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "the account email to login")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "the account password to login")
	rootCmd.AddCommand(logoutCmd)
}

