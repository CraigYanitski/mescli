package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/CraigYanitski/mescli/internal/requests"
	"github.com/CraigYanitski/mescli/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)



var getMessagesCmd = &cobra.Command{
    Use:   "messages",
    Short: "Get messages from the server",
    Long:  `Get messages stored on the server.

    This requires the user to be logged in.
    Any messages sent to the user are stored on the database until 
    retrieved. The messages will be displayed by user then by time.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("This command is not yet implemented...")
        messages, err := requests.GetMessages()
        if err != nil {
            return err
        }
        var u string
        for _, m := range messages {
            if m.SenderID.String() != u {
                u = m.SenderID.String()
                fmt.Printf("%s\n", utils.SuccessStyle.Bold(true).Render(u))
            }
            apiCfg.Messages[u] = append(
                apiCfg.Messages[u], 
                utils.RawMessage{
                    Sender: utils.ContactType,
                    Message: m.Message,
                    Time: m.CreatedAt,
                },
            )
            fmt.Printf(
                "  %s  %s\n", 
                utils.StatusStyle.Render(m.CreatedAt.Format("02-01-2006 15:04:05")), 
                m.Message,
            )
        }
        if ok := requests.WriteMessages(apiCfg.Messages); !ok {
            return errors.New("Unable to save messages locally.")
        }
        return nil
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
        fmt.Println("This command is not yet completed...")
        if len(apiCfg.Messages) == 0 {
            fmt.Println("There are no messages to display")
            return
        }
        for contact, messages := range apiCfg.Messages {
            fmt.Printf("%s\n", utils.SuccessStyle.Bold(true).Render(contact))
            for _, m := range messages {
                fmt.Printf(
                    "  %s  %s\n", 
                    utils.StatusStyle.Render(m.Time.Format("02-01-2006 15:04:05")), 
                    m.Message,
                )
            }
        }
    },
}

var sendMessageCmd = &cobra.Command{
    Use:   "message [MSG]",
    Short: "Send message to contact",
    Long:  `Send a message to a contact.

    The user email or uuid must be specified in order to send a message.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return fmt.Errorf("Require 1 argument, got %d", len(args))
        }
        msg := args[0]
        var uid *uuid.UUID
        if user == "" {
            return errors.New("must provide a user in order to send a message")
        } else if u, err := uuid.Parse(user); err == nil {
            uid = &u
            user, err = requests.GetContactEmail(u)
            if err != nil {
                return err
            }
        } else {
            uid, err = requests.GetContactID(user)
            if err != nil {
                return err
            }
        }
        packet, err := requests.AddContact(user)
        if err != nil {
            return err
        }
        err = requests.SendMessage(*uid, packet, msg)
        if err != nil {
            return err
        }
        apiCfg.Messages[uid.String()] = append(
            apiCfg.Messages[uid.String()], 
            utils.RawMessage{
                Sender: utils.SelfType,
                Message: msg,
                Time: time.Now(),
            },
        )
        if ok := requests.WriteMessages(apiCfg.Messages); !ok {
            return errors.New("Unable to save messages locally.")
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(sendMessageCmd)
    getCmd.AddCommand(getMessagesCmd)
    viewCmd.AddCommand(viewMessagesCmd)

    // Command flags
    rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user UUID or email")
}

