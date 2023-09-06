package db

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func HashPassword(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", fmt.Errorf("generate salt: %v", err)
	}

	return hashPasswordWithSalt(password, salt), nil
}

func ComparePassword(passwordHash, password string) (bool, error) {
	parts := strings.Split(passwordHash, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid password hash")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("decode salt: %v", err)
	}

	hash := hashPasswordWithSalt(password, salt)
	return hash == passwordHash, nil
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

func hashPasswordWithSalt(password string, salt []byte) string {
	const (
		time    = 1
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32
	)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	return fmt.Sprintf("%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
}
