package client

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"encoding/hex"

	//"fmt"
	"log"

	crypt "github.com/CraigYanitski/mescli/internal/cryptography"
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
    idkBytes := crypt.EncodeECDSAPublicKey(c.IdentityECDSA())

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
    idkBytes := crypt.EncodeECDSAPublicKey(c.IdentityECDSA())

    // encode ephemeral key in DER format
    epkBytes := crypt.EncodeECDHPublicKey(c.EphemeralKey())

    // return stringified keys
    return &MessagePacketJSON{
        IdentityKey: idkBytes, 
        EphemeralKey: epkBytes,
    }, nil
}

func ParsePrekeyPacket(packet *PrekeyPacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey, []byte, *ecdh.PublicKey) {
    rIKdsa := crypt.DecodeECDSAPublicKey(packet.IdentityKey)
    rSPK := crypt.DecodeECDHPublicKey(packet.SignedPrekey)
    rOK := crypt.DecodeECDHPublicKey(packet.OnetimePrekey)
    rSK, err := hex.DecodeString(packet.SignedKey)
    if err != nil {
        log.Fatal(err)
    }
    return rIKdsa, rSPK, rSK, rOK
}

func ParseMessagePacket(packet *MessagePacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey) {
    rIKdsa := crypt.DecodeECDSAPublicKey(packet.IdentityKey)
    rEPK := crypt.DecodeECDHPublicKey(packet.EphemeralKey)
    return rIKdsa, rEPK
}

