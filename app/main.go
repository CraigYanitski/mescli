package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CraigYanitski/mescli/internal/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func main() {
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

    // check if client is initialised
    //c := client.Client{
    //    Name: "test",//viper.GetString("name"),
    //}
    //if viper.Get("identity_token") == nil {
    //    c.Initialise(false)
    //}

    // bubble tea interface
    p := tea.NewProgram(InitialModel(), tea.WithAltScreen())

    // run
    m, err := p.Run()
    if err != nil {
        log.Fatal(err)
    }

    // run tests if selected
    if m, ok := m.(Model); ok && m.Chosen == 3 {
        runTests()
        fmt.Print("\n\n")
    }
}

// test encryption functionality
func runTests() {
    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("X3DH test")
    fmt.Printf("---------------\n")

    // Initialise clients in conversation
    alice := &client.Client{Name: "Alice"}
    _ = alice.Initialise(false)
    bob := &client.Client{Name: "Bob"}
    _ = bob.Initialise(true)
    log.Println("initialised")

    // get Bob's prekey package
    bobPKP, err := bob.SendPrekeyPacketJSON()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("have prekey packet")

    // Perform extended triple Diffie-Hellman exchange
    aliceMP := alice.InitiateX3DH(bobPKP, uuid.UUID{}, true)
    fmt.Printf("\nX3DH initialised\n")
    err = bob.CompleteX3DH(aliceMP, uuid.UUID{}, true)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nX3DH established\n")

    // Check if exchange was successful
    // Confirm whether or not they are equal, and thus the exchange is complete
    var result string
    if alice.CheckSecretEqual(bob) {
        result = successStyle.Italic(true).Render(
            "\nDiffie-Hellman secrets match - extended triple Diffie-Hellman exchange complete",
        )
    } else {
        result = errorStyle.Italic(true).Render(
            "\nDiffie-Hellman secrets don't match - error in establishing X3DH exchange! Secrets are not equal!!", 
        )
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Encryption test")
    fmt.Printf("---------------\n")

    fmt.Println("Alice -> Bob")

    // Try to send a message from Alice to Bob
    alicePub := alice.IdentityECDSA()
    bobPub := bob.IdentityECDSA()
    message := "Hi Bob!!"
    ciphertext, err := alice.SendMessage(message, bobPub, uuid.UUID{}, true)
    if err != nil {
        panic(err)
    }
    plaintext, err := bob.ReceiveMessage(ciphertext, alicePub, uuid.UUID{}, true)
    if err != nil {
        panic(err)
    }

    // Define progress strings
    initMessage := statusStyle.Bold(true).Render("\ninitial message (%d): ")
    initMessage += "%s\n"
    encrMessage := statusStyle.Bold(true).Render("\nencrypted message (%d): ")
    encrMessage += "0x%s\n"
    decrMessage := statusStyle.Bold(true).Render("\ndecrypted message (%d): ")
    decrMessage += "%s\n"

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result = successStyle.Italic(true).Render("\nMessage Encryption successful!!")
    } else {
        result = errorStyle.Italic(true).Render("\nError in message encryption!")
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Message length test")
    fmt.Printf("---------------\n")

    fmt.Println("Bob -> Alice")

    // Try to send a message from Alice to Bob
    message = "I am wondering about how much text I can put in a message before it encryption truncates. " +
              "There is obviously some entropy limit that cannot be surpassed given the SHA256 hashing function. " +
              "Perhaps this sentence will not make it through the transmission? " +
              "I should start splitting the message into chunks before finishing the encryption. " +
              "This message is clearly a good way to test this functionality."
    ciphertext, err = bob.SendMessage(message, alicePub, uuid.UUID{}, true)
    if err != nil {
        panic(err)
    }
    plaintext, err = alice.ReceiveMessage(ciphertext, bobPub, uuid.UUID{}, true)
    if err != nil {
        panic(err)
    }

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result = successStyle.Italic(true).Render("\nMessage Encryption successful!!")
    } else {
        result = errorStyle.Italic(true).Render("\nError in message encryption!")
    }
    fmt.Println(result)
}

