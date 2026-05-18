package cmd

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "os"

    "github.com/CraigYanitski/mescli/internal/tui"
    "github.com/CraigYanitski/mescli/internal/utils"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var verbose bool
var apiCfg *tui.ApiConfig

var user string
var email string
var password string

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
        // bubble tea interface
        p := tui.NewProgram(apiCfg)

        // run
        m, err := p.Run()
        if err != nil {
            log.Fatal(err)
        }

        // run tests if selected
        if m, ok := m.(tui.Model); ok && m.Chosen == 3 {
            utils.RunTests()
            fmt.Print("\n\n")
        } else {
			os.Exit(1)
		}
    },
}

func Execute() error {
    err := rootCmd.Execute()
    if err != nil {
        return errors.New("failed to run cobra command")
    }
    return nil
}

func init() {
    // set default user configuration
    //home, _ := os.UserHomeDir()
    //viper.AddConfigPath(path.Join(home, "projects/bootdev/courses/13-personal-project/mescli"))
    viper.AddConfigPath(".")
    viper.SetConfigName(".mescli")
    viper.SetConfigType("yaml")
    viper.SetDefault("api_url", "http://localhost:8080/api")
    viper.SetDefault("access_token", "")
    viper.SetDefault("refresh_token", "")
    viper.SetDefault("last_refresh", 0)
    viper.SetDefault("email", "")
    viper.SetDefault("name", "")
    viper.SetDefault("identity_key", "")
    viper.SetDefault("signed_prekey", "")
    viper.SetDefault("signed_key", "")
    //viper.SetDefault("root_ratchet", nil)
    //viper.SetDefault("send_ratchets", nil)
    //viper.SetDefault("recv_ratchets", nil)
    _ = viper.SafeWriteConfig()
    //if err != nil {
    //    log.Fatal(err)
    //}

    // load configuration
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalln(err)
    }
    viper.SetEnvPrefix("MESCLI_")
    viper.AutomaticEnv()

    apiCfg = &tui.ApiConfig{
        Name: viper.GetString("name"),
        Email: viper.GetString("email"),
    }

    messages := make(map[string][]utils.RawMessage)
    if _, err = os.Stat(".messages"); err == nil {
        msgBytes, err := os.ReadFile(".messages")
        if err != nil {
            log.Fatalln(err)
        }
        if err = json.Unmarshal(msgBytes, &messages); err != nil {
            log.Fatalln(err)
        }
        //fmt.Println(messages)
    }
    
    apiCfg.Messages = messages

    // check if client is initialised
    //c := client.Client{
    //    Name: "test",//viper.GetString("name"),
    //}
    //if viper.Get("identity_token") == nil {
    //    c.Initialise(false)
    //}


    // Command flags
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

