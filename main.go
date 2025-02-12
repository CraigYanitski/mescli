package main

import (
	"fmt"
	// "io"
	"log"
	"strings"

	"github.com/CraigYanitski/mescli/client"
	"github.com/CraigYanitski/mescli/typeset"
	// "github.com/charmbracelet/bubbles/list"
	// "github.com/charmbracelet/bubbles/textarea"
	// "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

func main() {
    // bubble tea interface
    p := tea.NewProgram(InitialModel())

    // run
    m, err := p.Run()
    if err != nil {
        log.Fatal(err)
    }

    // run tests if selected
    if m, ok := m.(Model); ok && m.Chosen == 2 {
        runTests()
    }

    // output additional padding
    fmt.Print("\n\n")
}

func runTests() {
    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("X3DH test")
    fmt.Printf("---------------\n")

    // Initialise clients in conversation
    alice := &client.Client{Name: "Alice"}
    _ = alice.Initialise()
    bob := &client.Client{Name: "Bob"}
    _ = bob.Initialise()

    // get Bob's prekey package
    bobPKP := bob.GetPrekeyPacket()

    // Perform extended triple Diffie-Hellman exchange
    aliceMP := alice.InitiateX3DH(bobPKP)
    fmt.Printf("\nX3DH initialised\n")
    err := bob.CompleteX3DH(aliceMP)
    if err != nil {
        panic(err)
    }
    fmt.Printf("\nX3DH established\n")

    // Check if exchange was successful
    // Confirm whether or not they are equal, and thus the exchange is complete
    var result string
    if alice.CheckSecretEqual(bob) {
        result, _ = typeset.FormatString("\nDiffie-Hellman secrets match - extended triple Diffie-Hellman exchange complete", 
            []string{"italics", "green"})
    } else {
        result, _ = typeset.FormatString("\nDiffie-Hellman secrets don't match - error in establishing X3DH exchange! Secrets are not equal!!", 
            []string{"italics", "red"})
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Encryption test")
    fmt.Printf("---------------\n")

    // Try to send a message from Alice to Bob
    alicePub, _ := alice.IdentityECDSA().ECDH()
    bobPub, _ := bob.IdentityECDSA().ECDH()
    message := "Hi Bob!!"
    ciphertext, err := alice.SendMessage(message, []string{"blue"}, bobPub)
    if err != nil {
        panic(err)
    }
    plaintext, err := bob.ReceiveMessage(ciphertext, alicePub)
    if err != nil {
        panic(err)
    }

    // Define progress strings
    initMessage, _ := typeset.FormatString("\ninitial message (%d): ", []string{"yellow", "bold"})
    initMessage += "%s\n"
    encrMessage, _ := typeset.FormatString("\nencrypted message (%d): ", []string{"yellow", "bold"})
    encrMessage += "0x%x\n"
    decrMessage, _ := typeset.FormatString("\ndecrypted message (%d): ", []string{"yellow", "bold"})
    decrMessage += "%s\n"

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result, _ = typeset.FormatString("\nMessage Encryption successful!!", []string{"green"})
    } else {
        result, _ = typeset.FormatString("\nError in message encryption!", []string{"red"})
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Message length test")
    fmt.Printf("---------------\n")

    // Try to send a message from Alice to Bob
    message = "I am wondering about how much text I can put in a message before it encryption truncates. " +
              "There is obviously some entropy limit that cannot be surpassed given the SHA256 hashing function. " +
              "Perhaps this sentence will not make it through the transmission? " +
              "I should start splitting the message into chunks before finishing the encryption. " +
              "This message is clearly a good way to test this functionality."
    ciphertext, err = alice.SendMessage(message, []string{}, bobPub)
    if err != nil {
        panic(err)
    }
    plaintext, err = bob.ReceiveMessage(ciphertext, alicePub)
    if err != nil {
        panic(err)
    }

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result, _ = typeset.FormatString("\nMessage Encryption successful!!", []string{"italics", "green"})
    } else {
        result, _ = typeset.FormatString("\nError in message encryption!", []string{"italics", "red"})
    }
    fmt.Println(result)
}

