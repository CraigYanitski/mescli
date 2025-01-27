package cryptography

import (
    "crypto/ecdh"
    "crypto/rand"
    "fmt"
    "golang.org/x/crypto/bcrypt"
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

