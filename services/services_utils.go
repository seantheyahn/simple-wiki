package services

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//GetSHA256Hex returns the sha256 hash of a []byte in hex format
func GetSHA256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//HashPassword makes the hash of the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash makes the hash from the password and compares to the input password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
