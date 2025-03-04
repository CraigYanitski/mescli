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
    //idk := c.IdentityECDSA()
    //idkBytes, err := x509.MarshalPKIXPublicKey(idk)
    //if err != nil {
    //    return nil, err
    //}
    idkBytes, err := crypt.SerialiseECDSAPublicKey(c.IdentityECDSA())
    if err != nil {
        return nil, err
    }

    // encode signed prekey in DER format
    spkBytes := crypt.SerialiseECDHPublicKey(c.SignedPrekey())
    //spkBytes, err := x509.MarshalPKIXPublicKey(&spk)
    //if err != nil {
    //    return nil, err
    //}

    // encode signed prekey in DER format
    opkBytes := crypt.SerialiseECDHPublicKey(c.OnetimePrekey())
    //opkBytes, err := x509.MarshalPKIXPublicKey(&opk)
    //if err != nil {
    //    return nil, err
    //}

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
    //idk := c.IdentityECDSA()
    //idkBytes, err := x509.MarshalPKIXPublicKey(idk)
    //if err != nil {
    //    return nil, err
    //}
    idkBytes, err := crypt.SerialiseECDSAPublicKey(c.IdentityECDSA())
    if err != nil {
        return nil, err
    }

    // encode ephemeral key in DER format
    epkBytes := crypt.SerialiseECDHPublicKey(c.EphemeralKey())
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
    //ridkInterface, err := x509.ParsePKIXPublicKey(packet.IdentityKey)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //rIKdsa, ok := ridkInterface.(*ecdsa.PublicKey)
    //if !ok {
    //    log.Fatal("error recasting identity key as ECDSA")
    //}
    rIKdsa, err := crypt.RecoverECDSAPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }

    //rSPK := contact.SignedPrekey
    //rspkInterface, err := x509.ParsePKIXPublicKey(packet.SignedPrekey)
    rSPK, err := crypt.RecoverECDHPublicKey(packet.SignedPrekey)
    if err != nil {
        log.Fatal(err)
    }
    //rSPK, ok := rspkInterface.(*ecdh.PublicKey)
    //if !ok {
    //    log.Fatal("error recasting signed prekey as ECDH")
    //}

    rSK, err := hex.DecodeString(packet.SignedKey)
    if err != nil {
        log.Fatal(err)
    }
    
    // rOK := contact.OnetimePrekey
    //rotkInterface, err := x509.ParsePKIXPublicKey(packet.OnetimePrekey)
    rOK, err := crypt.RecoverECDHPublicKey(packet.OnetimePrekey)
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
    //ridkInterface, err := x509.ParsePKIXPublicKey(packet.IdentityKey)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //rIKdsa, ok := ridkInterface.(*ecdsa.PublicKey)
    //if !ok {
    //    log.Fatal(err)
    //}
    rIKdsa, err := crypt.RecoverECDSAPublicKey(packet.IdentityKey)
    if err != nil {
        log.Fatal(err)
    }

    // rSPK := contact.SignedPrekey
    //repkInterface, err := x509.ParsePKIXPublicKey(packet.EphemeralKey)
    rEPK, err := crypt.RecoverECDHPublicKey(packet.EphemeralKey)
    if err != nil {
        log.Fatal(err)
    }
    //rEPK, ok := repkInterface.(*ecdh.PublicKey)
    //if !ok {
    //    log.Fatal(err)
    //}

    return rIKdsa, rEPK
}

//func SerialiseECDSAPublicKey(key *ecdsa.PublicKey) (string, error) {
//    keyBytes, err := x509.MarshalPKIXPublicKey(key)
//    if err != nil {
//        return "", err
//    }
//    return hex.EncodeToString(keyBytes), nil
//}
//
//func SerialiseECDSAPrivateKey(key *ecdsa.PrivateKey) (string, error) {
//    keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
//    if err != nil {
//        return "", err
//    }
//    return hex.EncodeToString(keyBytes), nil
//}
//
//func RecoverECDSAPublicKey(code string) (*ecdsa.PublicKey, error) {
//    keyBytes, err := hex.DecodeString(code)
//    if err != nil {
//        return nil, err
//    }
//    keyInterface, err := x509.ParsePKIXPublicKey(keyBytes)
//    if err != nil {
//        return nil, err
//    }
//    key, ok := keyInterface.(*ecdsa.PublicKey)
//    if !ok {
//        return nil, fmt.Errorf("error recasting bytes to ecdsa.PublicKey")
//    }
//    return key, nil
//}
//
//func RecoverECDSAPrivateKey(code string) (*ecdsa.PrivateKey, error) {
//    keyBytes, err := hex.DecodeString(code)
//    if err != nil {
//        return nil, err
//    }
//    keyInterface, err := x509.ParsePKCS8PrivateKey(keyBytes)
//    if err != nil {
//        return nil, err
//    }
//    key, ok := keyInterface.(*ecdsa.PrivateKey)
//    if !ok {
//        return nil, fmt.Errorf("error recasting bytes to ecdsa.PublicKey")
//    }
//    return key, nil
//}
//
//func SerialiseECDHPublicKey(key *ecdh.PublicKey) string {
//    return hex.EncodeToString(key.Bytes())
//}
//
//func SerialiseECDHPrivateKey(key *ecdh.PrivateKey) string {
//    return hex.EncodeToString(key.Bytes())
//}
//
//func RecoverECDHPublicKey(code string) (*ecdh.PublicKey, error) {
//    keyBytes, err := hex.DecodeString(code)
//    if err != nil {
//        return nil, err
//    }
//    key, err := ecdh.P256().NewPublicKey(keyBytes)
//    if err != nil {
//        return nil, err
//    }
//    return key, nil
//}
//
//func RecoverECDHPrivateKey(code string) (*ecdh.PrivateKey, error) {
//    keyBytes, err := hex.DecodeString(code)
//    if err != nil {
//        return nil, err
//    }
//    key, err := ecdh.P256().NewPrivateKey(keyBytes)
//    if err != nil {
//        return nil, err
//    }
//    return key, nil
//}

