package client

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"slices"

	crypt "github.com/CraigYanitski/mescli/cryptography"
	"github.com/CraigYanitski/mescli/typeset"
	"github.com/google/uuid"
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
    rootRatchet   *crypt.Ratchet
    sendRatchet   *crypt.Ratchet
    recvRatchet   *crypt.Ratchet
}

func (c *Client) Initialise(test bool) error {
    // check if keys are already in config or need to be generated
    if ik := viper.GetString("identity_key"); ik == "" || test {
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
            // save keys
            ikString := crypt.EncodeECDSAPrivateKey(ik)
            viper.Set("identity_key", ikString)
            spkString := crypt.EncodeECDHPrivateKey(spk)
            viper.Set("signed_prekey", spkString)
            viper.Set("signed_key", hex.EncodeToString(sk))
            opkString := crypt.EncodeECDHPrivateKey(opk)
            viper.Set("onetime_prekey", opkString)
            // write to config
            err = viper.WriteConfig()
            if err != nil {
                return fmt.Errorf("error saving cryptographic keys: %s", err)
            }
        }
    } else {
        // read keys from config
        c.identityKey = crypt.DecodeECDSAPrivateKey(viper.GetString("identity_key"))
        c.signedPrekey = crypt.DecodeECDHPrivateKey(viper.GetString("signed_prekey"))
        c.onetimePrekey = crypt.DecodeECDHPrivateKey(viper.GetString("onetime_prekey"))
        c.SignedKey, _ = hex.DecodeString(viper.GetString("signed_key"))
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
    //var key interface{} = c.signedPrekey.PublicKey()
    //if _, ok := key.(*ecdh.PublicKey); ok {
    //    log.Println("prekey type *ecdh.PublicKey")
    //}
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

func (c *Client) InitiateX3DH(contact *PrekeyPacketJSON, contactID uuid.UUID, test bool) *MessagePacketJSON {
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

    // check if ratchets are in config
    if rr := viper.GetString("root_ratchet"); rr == "" || test {
        // initialise root ratchet
        c.rootRatchet = &crypt.Ratchet{}
        c.rootRatchet.NewKDF(secret, nil, nil)

        // initialise sending ratchet
        sendSecret, _, err := c.rootRatchet.Extract(nil, nil, nil)
        if err != nil {
            log.Fatal(err)
        }
        c.sendRatchet = &crypt.Ratchet{}
        c.sendRatchet.NewKDF(sendSecret, nil, nil)

        // initialise receiving ratchet
        recvSecret, _, err := c.rootRatchet.Extract(nil, nil, nil)
        if err != nil {
            log.Fatal(err)
        }
        c.recvRatchet = &crypt.Ratchet{}
        c.recvRatchet.NewKDF(recvSecret, nil, nil)
    } else {
        // Make salt
        salt := make([]byte, len(secret))
        for i := 0; i < len(secret); i++ {
            salt[i] = secret[i] ^ 0xff
        }
        // get ratchets from config
        c.rootRatchet = crypt.DecodeRatchet(viper.GetString("contacts."+contactID.String()+".root_ratchet"), nil, nil)
        c.sendRatchet = crypt.DecodeRatchet(viper.GetString("contacts."+contactID.String()+".send_ratchet"), salt, nil)
        c.recvRatchet = crypt.DecodeRatchet(viper.GetString("contacts."+contactID.String()+".recv_ratchet"), salt, nil)
    }

    packet, err := c.SendMessagePacketJSON()
    if err != nil {
        log.Fatal(err)
    }
    
    // save ratchets in config
    if !test {
        // save secret 
        viper.Set("contacts."+contactID.String()+".secret", c.secret)
        // save ratchets
        viper.Set("contacts."+contactID.String()+".root_ratchet", c.rootRatchet.EncodeRatchet())
        viper.Set("contacts."+contactID.String()+".send_ratchets", c.sendRatchet.EncodeRatchet())
        viper.Set("contacts."+contactID.String()+".recv_ratchets", c.recvRatchet.EncodeRatchet())
        // write config
        err = viper.WriteConfig()
        if err != nil {
            log.Fatal(err)
        }
    }
    return packet
}

func (c *Client) CompleteX3DH(contact *MessagePacketJSON, contactID uuid.UUID, test bool) error {
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
    viper.Set("contacts."+contactID.String()+".secret", secret)

    // check if ratchets are in config
    if rr := viper.GetString("root_ratchet"); rr == "" || test {
        // initialise root ratchet
        c.rootRatchet = &crypt.Ratchet{}
        c.rootRatchet.NewKDF(secret, nil, nil)

        // initialise receiving ratchet
        recvSecret, _, err := c.rootRatchet.Extract(nil, nil, nil)
        if err != nil {
            return err
        }
        c.recvRatchet = &crypt.Ratchet{}
        c.recvRatchet.NewKDF(recvSecret, nil, nil)

        // initialise sending ratchet
        sendSecret, _, err := c.rootRatchet.Extract(nil, nil, nil)
        if err != nil {
            return err
        }
        c.sendRatchet = &crypt.Ratchet{}
        c.sendRatchet.NewKDF(sendSecret, nil, nil)
    } else {
        // Make salt
        salt := make([]byte, len(secret))
        for i := 0; i < len(secret); i++ {
            salt[i] = secret[i] ^ 0xff
        }
        // get ratchets from config
        c.rootRatchet = crypt.DecodeRatchet(viper.GetString("root_ratchet"), nil, nil)
        c.sendRatchet = crypt.DecodeRatchet(viper.GetString("send_ratchet"), salt, nil)
        c.recvRatchet = crypt.DecodeRatchet(viper.GetString("recv_ratchet"), salt, nil)
    }
    
    // save ratchets in config
    if !test {
        viper.Set("contacts."+contactID.String()+".secret", c.secret)
        // save ratchets
        viper.Set("contacts."+contactID.String()+".root_ratchet" , c.rootRatchet.EncodeRatchet())
        viper.Set("contacts."+contactID.String()+".recv_ratchets", c.recvRatchet.EncodeRatchet())
        viper.Set("contacts."+contactID.String()+".send_ratchets", c.sendRatchet.EncodeRatchet())
        // write config
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
    hash, err := crypt.HashPassword(password)

    // Return error if failed, else save password
    if err != nil {
        return err
    }
    c.password = hash
    return nil
}

func (c *Client) CheckPassword(password string) bool {
    // Hash password and compare to saved hash, return result
    ok := crypt.CheckPasswordHash(password, c.password)
    return ok
}

func (c *Client) SendMessage(plaintext string, format []string, pubkey *ecdh.PublicKey, contactID uuid.UUID, test bool) (string, error) {
    // Format message
    formattedMessage, err := typeset.FormatString(plaintext, format)
    if err != nil {
        return "", err
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
        return "", err
    }

    // Make salt
    salt := make([]byte, len(secret))
    for i := 0; i < len(secret); i++ {
        salt[i] = secret[i] ^ 0xff
    }

    // Generate key and iv
    sendKey, iv, err := c.sendRatchet.Extract(secret, salt, nil)
    if err != nil {
        return "", err
    }

    // Encrypt message
    ciphertext, err := crypt.EncryptMessage(sendKey, []byte(formattedMessage), iv)
    if err != nil {
        err = fmt.Errorf("error encrypting message: %v", err)
        return "", err
    }

    // save sending key if not test
    if !test {
        // save ratchet
        viper.Set("send_ratchet."+contactID.String()+".key", c.sendRatchet.EncodeRatchet())
        // write config
        viper.WriteConfig()
    }

    return hex.EncodeToString(ciphertext), nil
}

func (c *Client) ReceiveMessage(ciphertext string, pubkey *ecdh.PublicKey, contactID uuid.UUID, test bool) (string, error) {
    // Renew Diffie-Hellman key
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
    recvKey, iv, err := c.recvRatchet.Extract(secret, salt, nil)
    if err != nil {
        return "", err
    }

    // Decrypt message
    ciphertextBytes, err := hex.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    plaintext, err := crypt.DecryptMessage(recvKey, ciphertextBytes, iv)
    if err != nil {
        err = fmt.Errorf("error decrypting message: %v", err)
        return "", err
    }

    // save sending key if not test
    if !test {
        // save ratchet
        viper.Set("recv_ratchet."+contactID.String()+".key", c.recvRatchet.EncodeRatchet())
        // write config
        viper.WriteConfig()
    }

    return string(plaintext), nil
}

func generateECDH() (*ecdh.PrivateKey, error) {
    // generate private key
    key, err := crypt.GenerateECDH()

    // Return error if failed, else save key
    if err != nil {
        return nil, err
    }
    return key, nil
}

func generateECDSA() (*ecdsa.PrivateKey, error) {
    // generate private key
    key, err := crypt.GenerateECDSA()

    // Return error if failed, else save key
    if err != nil {
        return nil, err
    }
    return key, nil
}

