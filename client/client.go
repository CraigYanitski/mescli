package client

import (
    "crypto/ecdh"
    // "fmt"
    "github.com/CraigYanitski/mescli/cryptography"
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

