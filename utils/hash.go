package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	// We encode it to base64 so it can be safely stored as text in the database
	return base64.StdEncoding.EncodeToString(salt), nil
}

// Using SHA-256
func HashPassword(password string, salt string) string {
	combined := password + salt

	hasher := sha256.New()
	hasher.Write([]byte(combined))

	// Hexadecimal string
	return hex.EncodeToString(hasher.Sum(nil))
}

func CheckPassword(attemptedPassword string, storedSalt string, storedHash string) bool {

	attemptedHash := HashPassword(attemptedPassword, storedSalt)

	return attemptedHash == storedHash
}
