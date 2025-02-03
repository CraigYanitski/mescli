package main

import (
    "bytes"
    "fmt"
    "github.com/CraigYanitski/mescli/client"
    "github.com/CraigYanitski/mescli/typeset"
)

func main() {
    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("typesetting test")
    fmt.Printf("---------------\n\n")

    // Test string
    line := "A test string!!"

    // desired format
    format := []string{"default", "blue"}

    // Formatted string
    formattedLine, err := typeset.FormatString(line, format)
    if err != nil {
        panic(err)
    }

    // Test printing to terminal
    fmt.Println(formattedLine)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("X3DH test")
    fmt.Printf("---------------\n\n")

    // Initialise clients in conversation
    alice := client.Client{Name: "Alice"}
    _ = alice.Initialise()
    bob := client.Client{Name: "Bob"}
    _ = bob.Initialise()

    // Perform extended triple Diffie-Hellman exchange
    alice.EstablishX3DH(bob)
    fmt.Printf("%v's secret key: %x\n", alice.Name, alice.Secret)
    bob.CompleteX3DH(alice)
    fmt.Printf("%v's secret key: %x\n", bob.Name, bob.Secret)

    // Check if exchange was successful
    // Confirm whether or not they are equal, and thus the exchange is complete
    var result string
    if !bytes.Equal(alice.Secret, bob.Secret) {
        result, _ = typeset.FormatString("\nError in establishing X3DH exchange! Secrets are not equal!!", 
            []string{"italics", "red"})
    } else {
        result, _ = typeset.FormatString("\nExtended triple Diffie-Hellman exchange complete", 
            []string{"italics", "green"})
    }
    fmt.Println(result)
}
