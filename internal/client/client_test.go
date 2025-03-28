package client_test

import (
	"crypto"
	"crypto/ecdh"
	"fmt"
	"testing"

	"github.com/CraigYanitski/mescli/client"
	"github.com/google/uuid"
)

func TestClientCreation(t *testing.T) {
    type testCase struct {
        name      string
        password  string
        expected  bool
    }

    checkPublic := func (i interface{}) bool {
        _, ok := i.(crypto.PublicKey)
        return ok
    }
    checkPrivate := func (i interface{}) bool {
        i, ok := i.(ecdh.PrivateKey)
        return ok
    }

    tests := []testCase{
        {"Alice", "somelongpasswordforalice", true},
        {"Bob", "s0m310ngpa55w0rd4b0b", true},
    }

    failCount := 0
    passCount := 0

    fmt.Println("\n\nTesting Client creation")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("creating client %v with password %v", test.name, test.password)

        c := &client.Client{Name: test.name}
        c.HashPassword(test.password)
        c.Initialise(true)

        ik := c.IdentityECDSA()
        //if err != nil {
        //    t.Errorf("error getting client %v's public identity key: %v", c.Name, err)
        //    continue
        //}
        spk := c.SignedPrekey()
        //if err != nil {
        //    t.Errorf("error getting client %v's public signed prekey: %v", c.Name, err)
        //    continue
        //}
        opk := c.OnetimePrekey()
        //if err != nil {
        //    t.Errorf("error getting client %v's public one-time prekey: %v", c.Name, err)
        //    continue
        //}
        ek := c.EphemeralKey()
        //if err != nil {
        //    t.Errorf("error getting client %v's public ephemeral key: %v", c.Name, err)
        //    continue
        //}
        
        if checkPrivate(ik) {
            t.Error("error: client private identity key exposed!")
            continue
        }
        if checkPrivate(spk) {
            t.Error("error: client private signed prekey exposed!")
            continue
        }
        if checkPrivate(opk) {
            t.Error("error: client private one-time prekey exposed!")
            continue
        }
        if checkPrivate(ek) {
            t.Error("error: client private ephemeral key exposed!")
            continue
        }
        
        if !checkPublic(ik) {
            t.Error("error: client public identity key wrong type")
            continue
        }
        if !checkPublic(spk) {
            t.Error("error: client public signed prekey wrong type")
            continue
        }
        if !checkPublic(opk) {
            t.Error("error: client public one-time prekey wrong type")
            continue
        }
        if !checkPublic(ek) {
            t.Error("error: client public ephemeral key wrong type")
            continue
        }

        if !c.CheckPassword(test.password) {
            failCount++
            t.Errorf(`
Passwords not equal...
Inputs:    name: %v, password: %q, tryPassword: %q
Expected:  %v
Actual:    %v
`, test.name, test.password, test.password, test.expected, false)
        } else {
            passCount++
            fmt.Printf(`
Passwords equal
Inputs:    name: %v, password: %q, tryPassword: %q
Expected:  %v
Actual:    %v
`, test.name, test.password, test.password, test.expected, true)
        }

        tryPassword := test.password + " "

        if c.CheckPassword(tryPassword) {
            failCount++
            t.Errorf(`
Passwords shouldn't be equal...
Inputs:    name: %v, password: %q, tryPassword: %q
Expected:  %v
Actual:    %v
`, test.name, test.password, tryPassword, false, true)
        } else {
            passCount++
            fmt.Printf(`
Passwords are not equal
Inputs:    name: %v, password: %q, tryPassword: %q
Expected:  %v
Actual:    %v
`, test.name, test.password, tryPassword, false, false)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n\n\n", passCount, failCount)
}

func TestX3DH(t *testing.T) {
    type testCase struct {
        clientOneName  string
        clientTwoName  string
        expected       bool
    }

    tests := []testCase{
        {"Alice", "Bob", true},
    }

    failCount := 0
    passCount := 0

    fmt.Println("\n\nTesting Extended triple-Diffie Hellman exchange")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Creating client %v\n", test.clientOneName)

        clientOne := &client.Client{Name: test.clientOneName}
        err := clientOne.Initialise(true)
        if err != nil {
            t.Errorf("error initialising client %s's keys: %v", clientOne.Name, err)
        }

        fmt.Printf("Creating client %v\n", test.clientTwoName)

        clientTwo := &client.Client{Name: test.clientTwoName}
        err = clientTwo.Initialise(true)
        if err != nil {
            t.Errorf("error initialising client %s's keys: %v", clientTwo.Name, err)
        }

        fmt.Printf("%v initiating X3DH exchange\n", clientOne.Name)

        clientTwoPacket, err := clientTwo.SendPrekeyPacketJSON()
        if err != nil {
            t.Errorf("error sending client 1 message packet: %v", err)
        }

        clientOnePacket := clientOne.InitiateX3DH(clientTwoPacket, uuid.UUID{}, true)
        if err != nil {
            t.Errorf("error for %v initiating X3DH with %v: %v", clientOne.Name, clientTwo.Name, err)
        }

        fmt.Printf("%v completing X3DH exchange\n", clientTwo.Name)

        err = clientTwo.CompleteX3DH(clientOnePacket, uuid.UUID{}, true)
        if err != nil {
            t.Errorf("error for %v completing X3DH with %v: %v", clientTwo.Name, clientOne.Name, err)
        }

        result := clientOne.CheckSecretEqual(clientTwo)

        if result != test.expected {
            failCount++
            t.Errorf(`
Inputs:    clientOneName: %v, clientTwoName: %v
Expected:  X3DH established: %v
Actual:    X3DH established: %v
`, test.clientOneName, test.clientTwoName, true, result)
        } else {
            passCount++
            fmt.Printf(`
Inputs:    clientOneName: %v, clientTwoName: %v
Expected:  X3DH established: %v
Actual:    X3DH established: %v
`, test.clientOneName, test.clientTwoName, true, result)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n\n\n", passCount, failCount)
}
