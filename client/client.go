package client

import (
	"crypto/ecdh"
	"crypto/sha256"
	"fmt"
	"slices"
	"github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/typeset"
	"golang.org/x/crypto/hkdf"
)

type Client struct {
    Name           string
    password       string
    identityKey    *ecdh.PrivateKey
    signedPrekey   *ecdh.PrivateKey
    onetimePrekey  *ecdh.PrivateKey
    ephemeralKey   *ecdh.PrivateKey
    Secret         []byte
}

func (c *Client) Initialise() error {
    // generate identity key
    key, err := cryptography.GenerateECDH()
    if err != nil {
        return err
    }
    c.identityKey = key

    // generate signedPrekey
    key, err = cryptography.GenerateECDH()
    if err != nil {
        return err
    }
    c.signedPrekey = key

    // generate onetime prekey
    key, err = cryptography.GenerateECDH()
    if err != nil {
        return err
    }
    c.onetimePrekey = key

    // generate ephemeral key
    key, err = cryptography.GenerateECDH()
    if err != nil {
        return err
    }
    c.ephemeralKey = key

    return nil
}

func (c *Client) Identity() (*ecdh.PublicKey, error) {
    if c.identityKey == nil {
        return nil, fmt.Errorf("error returning identity key -- client not yet initialised")
    }
    return c.identityKey.PublicKey(), nil
}

func (c *Client) SignedPrekey() (*ecdh.PublicKey, error) {
    if c.signedPrekey == nil {
        return nil, fmt.Errorf("error returning signed prekey -- client not yet initialised")
    }
    return c.signedPrekey.PublicKey(), nil
}

func (c *Client) OnetimePrekey() (*ecdh.PublicKey, error) {
    if c.onetimePrekey == nil {
        return nil, fmt.Errorf("error returning one-time prekey -- client not yet initialised")
    }
    return c.onetimePrekey.PublicKey(), nil
}

func (c *Client) EphemeralKey() (*ecdh.PublicKey, error) {
    if c.ephemeralKey == nil {
        return nil, fmt.Errorf("error returning ephemeral key -- client not yet initialised")
    }
    return c.ephemeralKey.PublicKey(), nil
}

func (c *Client) EstablishX3DH(recipient Client) error {
    // get recipient public keys
    rIK, err := recipient.Identity()
    if err != nil {
        return err
    }
    rSK, err := recipient.SignedPrekey()
    if err != nil {
        return err
    }
    rOK, err := recipient.OnetimePrekey()
    if err != nil {
        return err
    }

    // verify signed prekey

    // calculate four DH secrets
    dh1, err := c.identityKey.ECDH(rSK)
    if err != nil {
        return err
    }
    dh2, err := c.ephemeralKey.ECDH(rIK)
    if err != nil {
        return err
    }
    dh3, err := c.ephemeralKey.ECDH(rSK)
    if err != nil {
        return err
    }
    dh4, err := c.ephemeralKey.ECDH(rOK)
    if err != nil {
        return err
    }

    // calculate secret key
    concat := slices.Concat(dh1, dh2, dh3, dh4)
    secret := make([]byte, 32)
    _, err = hkdf.New(sha256.New, concat, nil, nil).Read(secret)
    if err != nil {
        return err
    }

    // save secret and return
    c.Secret = secret
    return nil
}

func (c *Client) CompleteX3DH(sender Client) error {
    // get sender public keys
    sIK, err := sender.Identity()
    if err != nil {
        return err
    }
    sEK, err := sender.EphemeralKey()
    if err != nil {
        return err
    }

    // verify signed prekey

    // calculate four DH secrets
    dh1, err := c.signedPrekey.ECDH(sIK)
    if err != nil {
        return err
    }
    dh2, err := c.identityKey.ECDH(sEK)
    if err != nil {
        return err
    }
    dh3, err := c.signedPrekey.ECDH(sEK)
    if err != nil {
        return err
    }
    dh4, err := c.onetimePrekey.ECDH(sEK)
    if err != nil {
        return err
    }

    // calculate secret key
    concat := slices.Concat(dh1, dh2, dh3, dh4)
    secret := make([]byte, 32)
    _, err = hkdf.New(sha256.New, concat, nil, nil).Read(secret)
    if err != nil {
        return err
    }

    // save secret and return
    c.Secret = secret
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

