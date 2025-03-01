package client

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"slices"

	"github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/typeset"
	"github.com/spf13/viper"
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

func (c *Client) Initialise(test bool) error {
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

    // store keys in config file if not testing encryption
    if !test {
        //var send_ratchets map[string]cryptography.Ratchet
        //var recv_ratchets map[string]cryptography.Ratchet
        //viper.Set("identity_key", ik)
        //viper.Set("signed_prekey", spk)
        //viper.Set("signed_key", sk)
        //viper.Set("onetime_prekey", opk)
        //viper.Set("send_ratchets", send_ratchets)
        //viper.Set("recv_ratchets", recv_ratchets)
        err = viper.WriteConfig()
        if err != nil {
            return fmt.Errorf("error saving cryptographic keys: %s", err)
        }
    }

    return nil
}

func (c *Client) identityECDH() (*ecdh.PrivateKey) {
    if c.identityKey == nil {
        panic(fmt.Errorf("error returning identity key -- client not yet initialised"))
    }
    key, err := c.identityKey.ECDH()
    if err != nil {
        panic(err)
    }
    return key
}

func (c *Client) IdentityECDSA() (*ecdsa.PublicKey) {
    if c.identityKey == nil {
        panic(fmt.Errorf("error returning identity key -- client not yet initialised"))
    }
    return &c.identityKey.PublicKey
}

func (c *Client) SignedPrekey() (*ecdh.PublicKey) {
    if c.signedPrekey == nil {
        panic(fmt.Errorf("error returning signed prekey -- client not yet initialised"))
    }
    var key interface{} = c.signedPrekey.PublicKey()
    if _, ok := key.(*ecdh.PublicKey); ok {
        log.Println("prekey type *ecdh.PublicKey")
    }
    return c.signedPrekey.PublicKey()
}

func (c *Client) OnetimePrekey() (*ecdh.PublicKey) {
    if c.onetimePrekey == nil {
        panic(fmt.Errorf("error returning one-time prekey -- client not yet initialised"))
    }
    return c.onetimePrekey.PublicKey()
}

func (c *Client) EphemeralKey() (*ecdh.PublicKey) {
    if c.ephemeralKey == nil {
        panic(fmt.Errorf("error returning ephemeral key -- client not yet initialised"))
    }
    return c.ephemeralKey.PublicKey()
}

func (c *Client) InitiateX3DH(contact *PrekeyPacketJSON, test bool) *MessagePacketJSON {
    // get recipient identity public keys
    rIKdsa, rSPK, rSK, rOK := ParsePrekeyPacket(contact)
    rIK, err := rIKdsa.ECDH()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("parsed prekey packet")

    // verify signed prekey
    if !ecdsa.VerifyASN1(rIKdsa, encodeKey(rSPK), rSK) {
        log.Fatal("error verifying signed key during X3DH")
    }

    // generate ephemeral key
    ek, err := generateECDH()
    if err != nil {
        log.Fatal(err)
    }
    c.ephemeralKey = ek

    // get private ECDH
    iK := c.identityECDH()

    // calculate four DH secrets
    dh1, err := iK.ECDH(rSPK)
    if err != nil {
        log.Fatal(err)
    }
    dh2, err := c.ephemeralKey.ECDH(rIK)
    if err != nil {
        log.Fatal(err)
    }
    dh3, err := c.ephemeralKey.ECDH(rSPK)
    if err != nil {
        log.Fatal(err)
    }
    dh4, err := c.ephemeralKey.ECDH(rOK)
    if err != nil {
        log.Fatal(err)
    }

    // calculate secret key
    concat := slices.Concat(dh1, dh2, dh3, dh4)
    secret := make([]byte, 32)
    _, err = hkdf.New(sha256.New, concat, nil, nil).Read(secret)
    if err != nil {
        log.Fatal(err)
    }

    // save secret
    c.secret = secret

    // initialise root ratchet
    c.root_ratchet = &cryptography.Ratchet{}
    c.root_ratchet.NewKDF(secret, nil, nil)

    // initialise sending ratchet
    sendSecret, _, err := c.root_ratchet.Extract(nil, nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    c.send_ratchet = &cryptography.Ratchet{}
    c.send_ratchet.NewKDF(sendSecret, nil, nil)

    packet, err := c.SendMessagePacketJSON()
    if err != nil {
        log.Fatal(err)
    }
    
    // save ratchets in config
    if !test {
        //viper.Set("root_ratchet", c.root_ratchet)
        //send_ratchets := viper.GetStringMap("send_ratchets")
        //send_ratchets["contact"] = c.send_ratchet
        //viper.Set("send_ratchets", send_ratchets)
        err = viper.WriteConfig()
        if err != nil {
            log.Fatal(err)
        }
    }
    return packet
}

func (c *Client) CompleteX3DH(contact *MessagePacketJSON, test bool) error {
    // get sender public keys
    sIKdsa, sEK := ParseMessagePacket(contact)
    sIK, err := sIKdsa.ECDH()
    if err != nil {
        log.Fatal(err)
    }

    // get private ECDH key
    iK := c.identityECDH()

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
    
    // save ratchets in config
    if !test {
        //viper.Set("root_ratchet", c.root_ratchet)
        //recv_ratchets := viper.GetStringMap("recv_ratchets")
        //recv_ratchets["contact"] = c.send_ratchet
        //viper.Set("recv_ratchets", recv_ratchets)
        err = viper.WriteConfig()
        if err != nil {
            log.Fatal(err)
        }
    }

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
    key := c.identityECDH()

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
    key := c.identityECDH()

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

