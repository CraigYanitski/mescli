package client

import (
	"crypto/ecdh"
	"crypto/ecdsa"
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

type MessagePacket struct {
    // long-term 
    Identity   *ecdsa.PublicKey
    Ephemeral  *ecdh.PublicKey
    Message    []byte
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

func (c Client) GetMessagePacket() (*MessagePacket) {
    ik := c.IdentityECDSA()
    ek := c.EphemeralKey()
    return &MessagePacket{
        Identity: ik,
        Ephemeral: ek,
    }
}

