package services

import "crypto/sha256"

import "fmt"

//GetSHA256Hex returns the sha256 hash of a []byte in hex format
func GetSHA256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
