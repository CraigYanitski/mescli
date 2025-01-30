package client

import (
	"crypto/ecdh"
	"fmt"
	"github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/typesetting"
)

type Client struct {
    Name      string
    password  string
    KeyDH     *ecdh.PrivateKey
}

func (c *Client) GenerateKey() error {
    // generate private key
    key, err := cryptography.GenerateECDH()

    // Return error if failed, else save key
    if err != nil {
        return err
    }
    c.KeyDH = key
    return nil
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

func (c *Client) SendMessage(plaintext string, pubkey ecdh.PublicKey) ([]byte, error) {
    // Renew Diffie-Hellman key for encryption
    err := c.GenerateKey()
    if err != nil {
        err = fmt.Errorf("error generating DH key to send message: %v", err)
        return nil, err
    }

    // Generate nonce
    nonce := cryptography.GenerateNonce(15)

    // Encrypt message
    ciphertext, err := cryptography.EncryptMessage(c.KeyDH.ECDH(pubkey), []byte(plaintext), nonce)
    if err != nil {
        err = fmt.Errorf("error encrypting message: %v", err)
        return nil, err
    }

    return ciphertext, nil
}

func (c *Client) ReceiveMessage(ciphertext []byte, pubkey ecdh.PublicKey) (string, error) {
    // Renew Diffie-Hellman key
    err := c.GenerateKey()
    if err != nil {
        err = fmt.Errorf("error generating DH key for decryption: %v", err)
        return "", err
    }

    // Generate nonce
    nonce := cryptography.GenerateNonce(15)

    // Decrypt message
    plaintext, err := cryptography.DecryptMessage(c.KeyDH.ECDH(pubkey), ciphertext, nonce)
    if err != nil {
        err = fmt.Errorf("error decrypting message: %v", err)
        return "", err
    }

    return string(plaintext), nil
}

