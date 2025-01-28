package main

import (
    "bytes"
    "fmt"
    "github.com/CraigYanitski/mescli/client"
    "github.com/CraigYanitski/mescli/typeset"
)

func main() {
    /*
    // Initialise ANSI codes for test
    var codes = []typeset.AnsiCMD{3, 31}
    var reset = []typeset.AnsiCMD{0}

    // Get ANSI strings
    prefix, err := typeset.FormatANSI(codes)
    if err != nil {
        panic(err)
    }
    suffix, err := typeset.FormatANSI(reset)
    if err != nil {
        panic(err)
    }
    */

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

    // Separate tests
    fmt.Printf("\n---------------\n\n")

    // Create client 1
    alice := client.Client{Name: "Alice"}
    
    err = alice.HashPassword("Alice's very secure password")
    if err != nil {
        err = fmt.Errorf("Cannot hash Alice's password!! : %v", err)
        panic(err)
    }

    err = alice.GenerateKey()
    if err != nil {
        err = fmt.Errorf("Cannot generate Alice's key : %v", err)
        panic(err)
    }

    alicePrivKey := alice.KeyDH.Bytes()
    fmt.Printf("\nAlice's private key: %x\n", alicePrivKey)
    alicePubKey := alice.KeyDH.PublicKey()
    fmt.Printf("Alice's public key: %x\n", alicePubKey.Bytes())

    // Create client 2
    bob := client.Client{Name: "Bob"}
    
    err = bob.HashPassword("Bob's equally secure password")
    if err != nil {
        err = fmt.Errorf("Cannot hash Bob's password!! : %v", err)
        panic(err)
    }

    err = bob.GenerateKey()
    if err != nil {
        err = fmt.Errorf("Cannot generate Bob's key : %v", err)
        panic(err)
    }

    bobPrivKey := bob.KeyDH.Bytes()
    fmt.Printf("\nBob's private key: %x\n", bobPrivKey)
    bobPubKey := bob.KeyDH.PublicKey()
    fmt.Printf("Bob's public key: %x\n", bobPubKey.Bytes())

    // Calculate both clients' secrets
    //Alice
    aliceSecret, err := alice.KeyDH.ECDH(bobPubKey)

    if err != nil {
        err = fmt.Errorf("Error getting Alice's secret : %v", err)
        panic(err)
    }

    fmt.Printf("\nAlice's secret (%d): %x\n", len(aliceSecret), aliceSecret)

    // Bob
    bobSecret, err := bob.KeyDH.ECDH(alicePubKey)

    if err != nil {
        err = fmt.Errorf("Error getting Bob's secret : %v", err)
        panic(err)
    }
    fmt.Printf("Bob's secret (%d): %x\n", len(bobSecret), bobSecret)

    // Confirm whether or not they are equal, and thus the exchange is complete
    if !bytes.Equal(aliceSecret, bobSecret) {
        err = fmt.Errorf("Error in establishing DH exchange!! : %v", err)
        panic(err)
    }
    result, _ := typeset.FormatString("\nDiffie-Hellman exchange complete", []string{"italics", "green"})
    fmt.Println(result)

    // Renew keys in separate test
    fmt.Printf("\n---------------\n\n")

    err = alice.GenerateKey()
    if err != nil {
        err = fmt.Errorf("Cannot generate Alice's key : %v", err)
        panic(err)
    }

    alicePrivKey = alice.KeyDH.Bytes()
    fmt.Printf("\nAlice's private key: %x\n", alicePrivKey)
    alicePubKey = alice.KeyDH.PublicKey()
    fmt.Printf("Alice's public key: %x\n", alicePubKey.Bytes())

    err = bob.GenerateKey()
    if err != nil {
        err = fmt.Errorf("Cannot generate Bob's key : %v", err)
        panic(err)
    }

    bobPrivKey = bob.KeyDH.Bytes()
    fmt.Printf("\nBob's private key: %x\n", bobPrivKey)
    bobPubKey = bob.KeyDH.PublicKey()
    fmt.Printf("Bob's public key: %x\n", bobPubKey.Bytes())

    // Calculate both clients' secrets
    //Alice
    aliceSecret, err = alice.KeyDH.ECDH(bobPubKey)

    if err != nil {
        err = fmt.Errorf("Error getting Alice's secret : %v", err)
        panic(err)
    }

    fmt.Printf("\nAlice's secret (%d): %x\n", len(aliceSecret), aliceSecret)

    // Bob
    bobSecret, err = bob.KeyDH.ECDH(alicePubKey)

    if err != nil {
        err = fmt.Errorf("Error getting Bob's secret : %v", err)
        panic(err)
    }
    fmt.Printf("Bob's secret (%d): %x\n", len(bobSecret), bobSecret)

    // Confirm whether or not they are equal, and thus the exchange is complete
    if !bytes.Equal(aliceSecret, bobSecret) {
        err = fmt.Errorf("Error in establishing DH exchange!! : %v", err)
        panic(err)
    }
    result, _ = typeset.FormatString("\nDiffie-Hellman exchange complete", []string{"italics", "green"})
    fmt.Println(result)
}
