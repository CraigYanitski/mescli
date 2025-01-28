package cryptography

import (
    "crypto/ecdh"
    "crypto/rand"
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	l, err := rand.Read(salt)
	if (l != length) || (err != nil) {
		return nil, err
	}
	return salt[:], nil
}

func generateECDH() (*ecdh.PrivateKey, error) {
    // Generate private key
    key, err := ecdh.P256().GenerateKey(rand.Reader)

    // Return error if failed, else return key
    if err != nil {
        return nil, fmt.Errorf("error generating private key: %v", err)
    }
    return key, nil
}

