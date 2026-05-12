package utils

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CraigYanitski/mescli/internal/client"
	"github.com/google/uuid"
)

type SenderType int

const (
    SelfType SenderType = iota
    ContactType
)

type RawMessage struct {
    Sender   SenderType  `json:"sender"`
    Message  string      `json:"message"`
    Time     time.Time   `json:"time"`
}

// test encryption functionality
func RunTests() {
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
        result = SuccessStyle.Italic(true).Render(
            "\nDiffie-Hellman secrets match - extended triple Diffie-Hellman exchange complete",
        )
    } else {
        result = ErrorStyle.Italic(true).Render(
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
    initMessage := StatusStyle.Bold(true).Render("\ninitial message (%d): ")
    initMessage += "%s\n"
    encrMessage := StatusStyle.Bold(true).Render("\nencrypted message (%d): ")
    encrMessage += "0x%s\n"
    decrMessage := StatusStyle.Bold(true).Render("\ndecrypted message (%d): ")
    decrMessage += "%s\n"

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result = SuccessStyle.Italic(true).Render("\nMessage Encryption successful!!")
    } else {
        result = ErrorStyle.Italic(true).Render("\nError in message encryption!")
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
        result = SuccessStyle.Italic(true).Render("\nMessage Encryption successful!!")
    } else {
        result = ErrorStyle.Italic(true).Render("\nError in message encryption!")
    }
    fmt.Println(result)
}

func FormatOutput(raw string) string {
	paragraphs := strings.Split(raw, "\n\n")
	for i, p := range paragraphs {
		paragraphs[i] = strings.Join(strings.Split(p, "\n"), " ")
	}
	return strings.Join(paragraphs, "\n\n")
}

