package client

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
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
    identityKey    *ecdsa.PrivateKey
    signedPrekey   *ecdh.PrivateKey
    SignedKey      []byte
    onetimePrekey  *ecdh.PrivateKey
    ephemeralKey   *ecdh.PrivateKey
    secret         []byte
    root_ratchet   *cryptography.Ratchet
    send_ratchet   *cryptography.Ratchet
    recv_ratchet   *cryptography.Ratchet
}

func (c *Client) Initialise() error {
    // generate identity key
    ik, err := generateECDSA()
    if err != nil {
        return err
    }
    c.identityKey = ik

    // generate signedPrekey
    spk, err := generateECDH()
    if err != nil {
        return err
    }
    c.signedPrekey = spk

    // sign the prekey (required for sender verification)
    sk, err := ecdsa.SignASN1(rand.Reader, c.identityKey, encodeKey(c.signedPrekey.PublicKey()))
    if err != nil {
        return err
    }
    c.SignedKey = sk

    // generate onetime prekey
    opk, err := generateECDH()
    if err != nil {
        return err
    }
    c.onetimePrekey = opk

    return nil
}

func (c *Client) identityECDH() (*ecdh.PrivateKey, error) {
    if c.identityKey == nil {
        return nil, fmt.Errorf("error returning identity key -- client not yet initialised")
    }
    key, err := c.identityKey.ECDH()
    if err != nil {
        return nil, err
    }
    return key, nil
}

func (c *Client) IdentityECDH() (*ecdh.PublicKey, error) {
    if c.identityKey == nil {
        return nil, fmt.Errorf("error returning identity key -- client not yet initialised")
    }
    key, err := c.identityECDH()
    if err != nil {
        return nil, err
    }
    return key.PublicKey(), nil
}

func (c *Client) IdentityECDSA() (*ecdsa.PublicKey, error) {
    if c.identityKey == nil {
        return nil, fmt.Errorf("error returning identity key -- client not yet initialised")
    }
    // key, ok := c.identityKey.Public().(ecdsa.PublicKey)
    // if !ok {
    //     return nil, fmt.Errorf("error in recasting identity public key to ecdsa.PublicKey")
    // }
    return &c.identityKey.PublicKey, nil
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

func (c *Client) InitiateX3DH(recipient *Client) error {
    // get recipient public keys
    rIKdsa, err := recipient.IdentityECDSA()
    if err != nil {
        return err
    }
    rIK, err := recipient.IdentityECDH()
    if err != nil {
        return err
    }
    rSPK, err := recipient.SignedPrekey()
    if err != nil {
        return err
    }
    rSK := recipient.SignedKey
    rOK, err := recipient.OnetimePrekey()
    if err != nil {
        return err
    }

    // verify signed prekey
    if !ecdsa.VerifyASN1(rIKdsa, encodeKey(rSPK), rSK) {
        err = fmt.Errorf("error verifying signed key during X3DH")
        return err
    }

    // generate ephemeral key
    ek, err := generateECDH()
    if err != nil {
        return err
    }
    c.ephemeralKey = ek

    // get private ECDH
    iK, err := c.identityECDH()
    if err != nil {
        return err
    }

    // calculate four DH secrets
    dh1, err := iK.ECDH(rSPK)
    if err != nil {
        return err
    }
    dh2, err := c.ephemeralKey.ECDH(rIK)
    if err != nil {
        return err
    }
    dh3, err := c.ephemeralKey.ECDH(rSPK)
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

    // save secret
    c.secret = secret

    // initialise root ratchet
    c.root_ratchet = &cryptography.Ratchet{}
    c.root_ratchet.NewKDF(secret, nil, nil)

    // initialise sending ratchet
    sendSecret, _, err := c.root_ratchet.Extract(nil, nil, nil)
    if err != nil {
        return err
    }
    c.send_ratchet = &cryptography.Ratchet{}
    c.send_ratchet.NewKDF(sendSecret, nil, nil)
    return nil
}

func (c *Client) CompleteX3DH(sender *Client) error {
    // get sender public keys
    sIK, err := sender.IdentityECDH()
    if err != nil {
        return err
    }
    sEK, err := sender.EphemeralKey()
    if err != nil {
        return err
    }

    // get private ECDH key
    iK, err := c.identityECDH()
    if err != nil {
        return err
    }

    // calculate four DH secrets
    dh1, err := c.signedPrekey.ECDH(sIK)
    if err != nil {
        return err
    }
    dh2, err := iK.ECDH(sEK)
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
    c.secret = secret

    // initialise root ratchet
    c.root_ratchet = &cryptography.Ratchet{}
    c.root_ratchet.NewKDF(secret, nil, nil)

    // initialise receiving ratchet
    recvSecret, _, err := c.root_ratchet.Extract(nil, nil, nil)
    if err != nil {
        return err
    }
    c.recv_ratchet = &cryptography.Ratchet{}
    c.recv_ratchet.NewKDF(recvSecret, nil, nil)
    return nil
}

func (c *Client) CheckSecretEqual(contact *Client) bool {
    return bytes.Equal(c.secret, contact.secret)
}

func encodeKey(key *ecdh.PublicKey) []byte {
    h := sha256.New()
    h.Write(key.Bytes())
    return h.Sum(nil)
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
    // key, err := generateKey()
    // if err != nil {
    //     err = fmt.Errorf("error generating DH key to send message: %v", err)
    //     return nil, err
    // }
    key, err := c.identityECDH()
    if err != nil {
        return nil, err
    }

    // Calculate shared secret
    secret, err := key.ECDH(pubkey)
    if err != nil {
        return nil, err
    }

    // Make salt
    salt := make([]byte, len(secret))
    for i := 0; i < len(secret); i++ {
        salt[i] = secret[i] ^ 0xff
    }

    // Generate key and iv
    sendKey, iv, err := c.send_ratchet.Extract(secret, salt, nil)
    if err != nil {
        return nil, err
    }

    // Encrypt message
    ciphertext, err := cryptography.EncryptMessage(sendKey, []byte(formattedMessage), iv)
    if err != nil {
        err = fmt.Errorf("error encrypting message: %v", err)
        return nil, err
    }

    return ciphertext, nil
}

func (c *Client) ReceiveMessage(ciphertext []byte, pubkey *ecdh.PublicKey) (string, error) {
    // Renew Diffie-Hellman key
    // key, err := generateECDSA()
    // if err != nil {
    //     err = fmt.Errorf("error generating DH key for decryption: %v", err)
    //     return "", err
    // }
    key, err := c.identityECDH()
    if err != nil {
        return "", err
    }

    // Calculate shared secret
    secret, err := key.ECDH(pubkey)
    if err != nil {
        return "", err
    }

    // Make salt
    salt := make([]byte, len(secret))
    for i := 0; i < len(secret); i++ {
        salt[i] = secret[i] ^ 0xff
    }

    // Generate key and iv
    recvKey, iv, err := c.recv_ratchet.Extract(secret, salt, nil)
    if err != nil {
        return "", err
    }

    // Decrypt message
    plaintext, err := cryptography.DecryptMessage(recvKey, ciphertext, iv)
    if err != nil {
        err = fmt.Errorf("error decrypting message: %v", err)
        return "", err
    }

    return string(plaintext), nil
}

func generateECDH() (*ecdh.PrivateKey, error) {
    // generate private key
    key, err := cryptography.GenerateECDH()

    // Return error if failed, else save key
    if err != nil {
        return nil, err
    }
    return key, nil
}

func generateECDSA() (*ecdsa.PrivateKey, error) {
    // generate private key
    key, err := cryptography.GenerateECDSA()

    // Return error if failed, else save key
    if err != nil {
        return nil, err
    }
    return key, nil
}

