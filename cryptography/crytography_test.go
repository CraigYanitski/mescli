package cryptography_test

import (
	"crypto/ecdh"
	"fmt"
	"testing"

	"github.com/CraigYanitski/mescli/cryptography"
)

func TestPasswordHash(t *testing.T) {
    type testCase struct {
        password     string
        tryPassword  string
        expected     bool
    }

    tests := []testCase{}

    failCount := 0
    passCount := 0

    fmt.Println("Testing X")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Hashing password %q and comparing to %q", test.password, test.tryPassword)

        hash, err := cryptography.HashPassword(test.password)
        if err != nil {
            t.Errorf("error hashing password using bcrypt: %v", err)
            continue
        }

        fmt.Printf("Hashed password: %s", hash)

        result := cryptography.CheckPasswordHash(test.password, hash)

        if result != test.expected {
            failCount++
            t.Errorf(`
Inputs:    password: %q, tryPassword: %q
Expected:  %t
Actual:    %t`, test.password, test.tryPassword, test.expected, result)
        } else {
            passCount++
            fmt.Printf(`
Inputs:    password: %q, tryPassword: %q
Expected:  %t
Actual:    %t`, test.password, test.tryPassword, test.expected, result)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestDHKeyGeneration(t *testing.T) {
    type testCase struct {
        expected  bool
    }

    tests := []testCase{
        {true},
    }

    failCount := 0
    passCount := 0

    fmt.Println("Testing Diffie-Hellman key generation")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Generating Diffie-Hellman key")

        result, err := cryptography.GenerateECDH()
        if err != nil {
            t.Errorf("error generating DH key: %v", err)
        }

        var keyPrivate interface{} = result
        _, okPrivate := keyPrivate.(*ecdh.PrivateKey)
        var keyPublic interface{} = result.Public()
        _, okPublic := keyPublic.(*ecdh.PublicKey)

        if (okPrivate && okPublic) != test.expected {
            failCount++
            t.Errorf(`
Inputs:    %v
Expected:  %v
Actual:    %v`, nil, test.expected, result)
        } else {
            passCount++
            fmt.Printf(`
Inputs:    %v
Expected:  %v
Actual:    %v`, nil, test.expected, result)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestEncryptMessage(t *testing.T) {
    type testCase struct {
        message   string
        expected  string
    }

    tests := []testCase{
        {"Some short message", "random bytes"},
        {"Some very long message " +
         "that goes on and on and on." +
         " This should also trigger the ciphertext to get longer.", "some long excrypted data"},
    }

    failCount := 0
    passCount := 0

    fmt.Println("Testing message encryption")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Encrypting message %q", test.message)

        // create new ratchet to read encryption keys
        r := cryptography.Ratchet{}
        r.NewKDF(nil, nil, nil)
        key, iv, err := r.Extract(nil, nil, nil)
        if err != nil {
            t.Errorf("error extracting from KDF: %v", err)
            continue
        }

        // encrypt the message
        ciphertext, err := cryptography.EncryptMessage(key, []byte(test.message), iv)
        if err != nil {
            t.Errorf("error encrypting message: %v", err)
        }

        if string(ciphertext) != test.message {
            failCount++
            t.Errorf(`
Inputs:    %v
Expected:  %v
Actual:    %x`, test.message, test.expected, ciphertext)
        } else {
            passCount++
            fmt.Printf(`
Inputs:    %v
Expected:  %v
Actual:    %x`, test.message, test.expected, ciphertext)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestRatchetExtraction(t *testing.T) {
    type testCase struct {
        secret    []byte
        expected  []int
    }

    test := testCase{[]byte{}, []int{32, 15}}

    //failCount := 0
    //passCount := 0

    fmt.Println("Testing cryptographic ratchet creation and key extraction")

    //for _, test := range tests {
    fmt.Println("----------------------------------------")
    fmt.Printf("Generating ratchet with %x", test.secret)

    ratchet := cryptography.Ratchet{}
    ratchet.NewKDF(test.secret, nil, nil)
    key, iv, err := ratchet.Extract(nil, nil, nil)
    if err != nil {
        t.Errorf("error Reading key from KDF: %v", err)
    }

    result := (len(key) == test.expected[0]) || (len(iv) == test.expected[1])

    if !result {
        //failCount++
        t.Errorf(`
Inputs:    secret: %v
Expected:  %v
Actual:    %v`, test.secret, test.expected, []int{len(key), len(iv)})
    } else {
        //passCount++
        fmt.Printf(`
Inputs:    %v
Expected:  %v
Actual:    %v`, test.secret, test.expected, []int{len(key), len(iv)})
        //}
    }

    fmt.Println("========================================")
    //fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

