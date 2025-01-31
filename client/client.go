package client

import (
	"crypto/ecdh"
	"fmt"
	"github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/typeset"
)

type Client struct {
    Name      string
    password  string
    KeyDH     *ecdh.PrivateKey
}

func (c *Client) HashPassword(password string) error {
    // Hash password
    hash, err := cryptography.HashPassword(password)

    // Return error if failed, else save password
    if err != nil {
        return err
    }
    c.password = hash
    return nil
}

func (c *Client) CheckPassword(password string) bool {
    // Hash password and compare to saved hash, return result
    ok := cryptography.CheckPasswordHash(password, c.password)
    return ok
}

func (c *Client) SendMessage(plaintext string, format []string, pubkey *ecdh.PublicKey) ([]byte, error) {
    // Format message
    formattedMessage, err := typeset.FormatString(plaintext, format)
    if err != nil {
        return nil, err
    }

    // Renew Diffie-Hellman key for encryption
    key, err := generateKey()
    if err != nil {
        err = fmt.Errorf("error generating DH key to send message: %v", err)
        return nil, err
    }

    // Generate nonce
    nonce := cryptography.GenerateNonce(15)

    // Calculate shared secret
    secret, err := key.ECDH(pubkey)
    if err != nil {
        return nil, err
    }

    // Encrypt message
    ciphertext, err := cryptography.EncryptMessage(secret, []byte(formattedMessage), nonce)
    if err != nil {
        err = fmt.Errorf("error encrypting message: %v", err)
        return nil, err
    }

    return ciphertext, nil
}

func (c *Client) ReceiveMessage(ciphertext []byte, pubkey *ecdh.PublicKey) (string, error) {
    // Renew Diffie-Hellman key
    key, err := generateKey()
    if err != nil {
        err = fmt.Errorf("error generating DH key for decryption: %v", err)
        return "", err
    }

    // Generate nonce
    nonce := cryptography.GenerateNonce(15)

    // Calculate shared secret
    secret, err := key.ECDH(pubkey)
    if err != nil {
        return "", err
    }

    // Decrypt message
    plaintext, err := cryptography.DecryptMessage(secret, ciphertext, nonce)
    if err != nil {
        err = fmt.Errorf("error decrypting message: %v", err)
        return "", err
    }

    return string(plaintext), nil
}

func generateKey() (*ecdh.PrivateKey, error) {
    // generate private key
    key, err := cryptography.GenerateECDH()

    // Return error if failed, else save key
    if err != nil {
        return nil, err
    }
    return key, nil
}

