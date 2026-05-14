package cmd

import (
    "errors"
    "fmt"

    "github.com/CraigYanitski/mescli/internal/requests"
    "github.com/CraigYanitski/mescli/internal/utils"
    "github.com/spf13/cobra"
)



var getUserCmd = &cobra.Command{
    Use:   "user",
    Short: "Get user information from the server",
    Long:  `Get user information from the server.

    The user email must be specified.
    This prints the user's UUID.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("This command is not yet completed...")
        var email string
        if len(args) !=1 {
            return errors.New("A user email must be specified to get their UUID")
        } else {
            email = args[0]
        }
        u, err := requests.GetUser(email)
        if err != nil {
            return err
        }
        fmt.Printf("User %s (%s) \n", u.ID.String(), email)
        return nil
    },
}

var viewUsersCmd = &cobra.Command{
    Use:   "users",
    Short: "View user information from local disk",
    Long:  `View user information from local disk.

    The user email may be specified for specific information.
    This prints each user's UUID, message statistics, and most recent message.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("This command is not yet completed...")
        if len(apiCfg.Messages) == 0 {
            fmt.Println("There are no messages to display")
            return
        }
        for contact := range apiCfg.Messages {
            fmt.Printf("%s\n", utils.SuccessStyle.Bold(true).Render(contact))
        }
    },
}

func init() {
    getCmd.AddCommand(getUserCmd)
    viewCmd.AddCommand(viewUsersCmd)
}

