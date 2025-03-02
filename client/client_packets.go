package client

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"log"
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
    IdentityKey    []byte  `json:"identity_key"`
    SignedPrekey   []byte  `json:"signed_prekey"`
    SignedKey      []byte  `json:"signed_key"`
    OnetimePrekey  []byte  `json:"onetime_prekey"`
}

type MessagePacket struct {
    // long-term 
    Identity   *ecdsa.PublicKey
    Ephemeral  *ecdh.PublicKey
    Message    []byte
}

type MessagePacketJSON struct {
    IdentityKey   []byte  `json:"identity_key"`
    EphemeralKey  []byte  `json:"ephemeral_key"`
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
    idk := c.IdentityECDSA()
    idkBytes, err := x509.MarshalPKIXPublicKey(idk)
    if err != nil {
        return nil, err
    }

    // encode signed prekey in DER format
    spkBytes := c.SignedPrekey().Bytes()
    //spkBytes, err := x509.MarshalPKIXPublicKey(&spk)
    //if err != nil {
    //    return nil, err
    //}

    // encode signed prekey in DER format
    opkBytes := c.OnetimePrekey().Bytes()
    //opkBytes, err := x509.MarshalPKIXPublicKey(&opk)
    //if err != nil {
    //    return nil, err
    //}

    // encode signed key in DER format
    skBytes := c.SignedKey

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
    idk := c.IdentityECDSA()
    idkBytes, err := x509.MarshalPKIXPublicKey(idk)
    if err != nil {
        return nil, err
    }

    // encode ephemeral key in DER format
    epkBytes := c.EphemeralKey().Bytes()
    //epkBytes, err := x509.MarshalPKIXPublicKey(epk)
    //if err != nil {
    //    return nil, err
    //}

    // return stringified keys
    return &MessagePacketJSON{
        IdentityKey: idkBytes, 
        EphemeralKey: epkBytes,
    }, nil
}

func ParsePrekeyPacket(packet *PrekeyPacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey, []byte, *ecdh.PublicKey) {
    ridkInterface, err := x509.ParsePKIXPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }
    rIKdsa, ok := ridkInterface.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("error recasting identity key as ECDSA")
    }

    //rSPK := contact.SignedPrekey
    //rspkInterface, err := x509.ParsePKIXPublicKey(packet.SignedPrekey)
    rSPK, err := ecdh.P256().NewPublicKey(packet.SignedPrekey)
    if err != nil {
        log.Fatal(err)
    }
    //rSPK, ok := rspkInterface.(*ecdh.PublicKey)
    //if !ok {
    //    log.Fatal("error recasting signed prekey as ECDH")
    //}

    rSK := packet.SignedKey
    
    // rOK := contact.OnetimePrekey
    //rotkInterface, err := x509.ParsePKIXPublicKey(packet.OnetimePrekey)
    rOK, err := ecdh.P256().NewPublicKey(packet.OnetimePrekey)
    if err != nil {
        log.Fatal(err)
    }
    //rOK, ok := rotkInterface.(*ecdh.PublicKey)
    //if !ok {
    //    log.Fatal("error recasting onetime key as ECDH")
    //}

    return rIKdsa, rSPK, rSK, rOK
}

func ParseMessagePacket(packet *MessagePacketJSON) (*ecdsa.PublicKey, *ecdh.PublicKey) {
    ridkInterface, err := x509.ParsePKIXPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }
    rIKdsa, ok := ridkInterface.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal(err)
    }

    // rSPK := contact.SignedPrekey
    //repkInterface, err := x509.ParsePKIXPublicKey(packet.EphemeralKey)
    rEPK, err := ecdh.P256().NewPublicKey(packet.EphemeralKey)
    if err != nil {
        log.Fatal(err)
    }
    //rEPK, ok := repkInterface.(*ecdh.PublicKey)
    //if !ok {
    //    log.Fatal(err)
    //}

    return rIKdsa, rEPK
}

