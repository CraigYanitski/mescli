package client

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"encoding/hex"
	//"fmt"
	"log"

	crypt "github.com/CraigYanitski/mescli/cryptography"
)

type PrekeyPacket struct {
    // long-term DSA identity key
    Identity    *ecdsa.PublicKey
    // signed prekey used in signature
    SignedPrekey   *ecdh.PublicKey
    // onetime prekey used for X3DH encryption
    OnetimePrekey  *ecdh.PublicKey
    // signature essential for contact verification
    SignedKey      []byte
}

type PrekeyPacketJSON struct {
    IdentityKey    string  `json:"identity_key"`
    SignedPrekey   string  `json:"signed_prekey"`
    SignedKey      string  `json:"signed_key"`
    OnetimePrekey  string  `json:"onetime_prekey"`
}

type MessagePacket struct {
    // long-term 
    Identity   *ecdsa.PublicKey
    Ephemeral  *ecdh.PublicKey
    Message    []byte
}

type MessagePacketJSON struct {
    IdentityKey   string  `json:"identity_key"`
    EphemeralKey  string  `json:"ephemeral_key"`
}

func (c Client) GetPrekeyPacket() (*PrekeyPacket) {
    ik := c.IdentityECDSA()
    spk := c.SignedPrekey()
    opk := c.OnetimePrekey()
    return &PrekeyPacket {
        Identity: ik,
        SignedPrekey: spk,
        OnetimePrekey: opk,
        SignedKey: c.SignedKey,
    }
}

func (c Client) SendPrekeyPacketJSON() (*PrekeyPacketJSON, error) {
    // encode identity key in DER format
    idkBytes, err := crypt.EncodeECDSAPublicKey(c.IdentityECDSA())
    if err != nil {
        return nil, err
    }

    // encode signed prekey in DER format
    spkBytes := crypt.EncodeECDHPublicKey(c.SignedPrekey())

    // encode signed prekey in DER format
    opkBytes := crypt.EncodeECDHPublicKey(c.OnetimePrekey())

    // encode signed key in DER format
    skBytes := hex.EncodeToString(c.SignedKey)

    // return stringified keys
    return &PrekeyPacketJSON{
        IdentityKey: idkBytes, 
        SignedPrekey: spkBytes,
        SignedKey: skBytes,
        OnetimePrekey: opkBytes,
    }, nil
}

func (c Client) GetMessagePacket() (*MessagePacket) {
    ik := c.IdentityECDSA()
    ek := c.EphemeralKey()
    return &MessagePacket{
        Identity: ik,
        Ephemeral: ek,
    }
}

func (c Client) SendMessagePacketJSON() (*MessagePacketJSON, error) {
    // encode identity key in DER format
    idkBytes, err := crypt.EncodeECDSAPublicKey(c.IdentityECDSA())
    if err != nil {
        return nil, err
    }

    // encode ephemeral key in DER format
    epkBytes := crypt.EncodeECDHPublicKey(c.EphemeralKey())

    // return stringified keys
    return &MessagePacketJSON{
        IdentityKey: idkBytes, 
        EphemeralKey: epkBytes,
    }, nil
}

func ParsePrekeyPacket(packet *PrekeyPacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey, []byte, *ecdh.PublicKey) {
    rIKdsa, err := crypt.DecodeECDSAPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }

    rSPK, err := crypt.DecodeECDHPublicKey(packet.SignedPrekey)
    if err != nil {
        log.Fatal(err)
    }

    rSK, err := hex.DecodeString(packet.SignedKey)
    if err != nil {
        log.Fatal(err)
    }
    
    rOK, err := crypt.DecodeECDHPublicKey(packet.OnetimePrekey)
    if err != nil {
        log.Fatal(err)
    }

    return rIKdsa, rSPK, rSK, rOK
}

func ParseMessagePacket(packet *MessagePacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey) {
    rIKdsa, err := crypt.DecodeECDSAPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }

    rEPK, err := crypt.DecodeECDHPublicKey(packet.EphemeralKey)
    if err != nil {
        log.Fatal(err)
    }

    return rIKdsa, rEPK
}

