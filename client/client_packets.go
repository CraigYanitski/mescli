package client

import (
	"crypto/ecdh"
	"crypto/ecdsa"
)

type PreKeyPacket struct {
    // long-term identity key
    Identity       ecdsa.PublicKey
    // signed prekey used in signature
    SignedPrekey   ecdh.PublicKey
    // onetime prekey used for X3DH encryption
    OnetimePrekey  ecdh.PublicKey
    // signature essential for contact verification
    SignedKey      []byte
}

type MessagePacket struct {
    Identity   ecdsa.PublicKey
    Ephemeral  ecdh.PublicKey
    Message    []byte
}

