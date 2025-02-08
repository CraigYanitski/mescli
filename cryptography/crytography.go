package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/hkdf"
)

const (
    NonceSize = 15
)

func HashPassword(password string) (string, error) {
	// Hash password a specified number of times
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)

    // return error if failed, else return hash as a string
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	// Compare password to hash
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
    // Print error if failed and return false, else return true
    if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func generateSalt(length int) ([]byte, error) {
    // Initialise salt slice of specified length
	salt := make([]byte, length)

    // Generate random salt
	l, err := rand.Read(salt)

    // Return error if failed, else return salt
	if (l != length) || (err != nil) {
		return nil, err
	}
	return salt[:], nil
}

func GenerateECDH() (*ecdh.PrivateKey, error) {
    // Generate private key
    key, err := ecdh.P256().GenerateKey(rand.Reader)

    // Return error if failed, else return key
    if err != nil {
        return nil, fmt.Errorf("error generating private key: %v", err)
    }
    return key, nil
}

func GenerateECDSA() (*ecdsa.PrivateKey, error) {
    // Generate private key
    key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

    // Return error if failed, else return key
    if err != nil {
        return nil, fmt.Errorf("error generating private key: %v", err)
    }
    return key, nil
}

func GenerateNonce(length int) []byte {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return randomBytes
}

func EncryptMessage(key, plaintext, nonce []byte) (ciphertext []byte, err error) {
    // Create new cipher block
    block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

    // create new GCM cipher
	aesgcm, err := cipher.NewGCMWithNonceSize(block, NonceSize)
	if err != nil {
		return nil, err
	}

    // encrypt the plaintext
	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func DecryptMessage(key, ciphertext, nonce []byte) (plaintext []byte, err error) {
	// Create a new cipher block
    block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

    // create new GCM cipher
	aesaead, err := cipher.NewGCMWithNonceSize(block, NonceSize)
	if err != nil {
		return nil, err
	}

    // decrypt the ciphertext
	plaintext, err = aesaead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

type Ratchet struct {
    // This is a KDF reader that will be used to generate new keys
    kdf  io.Reader
    // This is the root key to be concatenated with the input data
    key  []byte
}

// This is a function to create a new KDF Reader from which keys can be read
func (r *Ratchet) NewKDF(secret, salt, info []byte) {
    r.key = secret
    r.kdf = hkdf.New(sha256.New, secret, salt, info)
}

// This function performs an extract and expand on the KDF to derive a new key and initialisation vector
func (r *Ratchet) Extract(input, salt, info []byte) (key []byte, iv []byte, err error) {
    secret := append(r.key, input...)
    // kdf := hkdf.Extract(sha256.New, secret, salt)
    kdfKey := hkdf.Extract(sha256.New, secret, salt)
    r.key = kdfKey
    r.kdf = hkdf.Expand(sha256.New, kdfKey, info)
    key = make([]byte, 32)
    iv = make([]byte, NonceSize)
    _, err = r.kdf.Read(key)
    if err != nil {
        return nil, nil, err
    }
    _, err = r.kdf.Read(iv)
    if err != nil {
        return nil, nil, err
    }
    return key, iv, nil
}

