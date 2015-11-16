package password

// This will handle all aspects of authenticating users in our system
// For password managing/salting I used:
// http://austingwalters.com/building-a-web-server-in-go-salting-passwords/

import (
	"crypto/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	SaltLength = 64
	// On a scale of 3 - 31, how intense Bcrypt should be
	EncryptCost = 14
)

// Hashes password using the salt, then tacks it onto the front of the string
// so we can reverse engineer it with a correct password
func hashPassword(salt string, password string) ([]byte, error) {
	combo := combine(salt, password)
	return bcrypt.GenerateFromPassword([]byte(combo), EncryptCost)
}

// Generates a random salt using DevNull
func generateSalt() (string, error) {

	// Read in data
	data := make([]byte, SaltLength)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func combine(salt string, hash string) string {
	return strings.Join([]string{salt, hash}, "")
}

// Encrypt creates a new hash/salt combo from a raw password as inputted
// by the user
func Encrypt(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash, err := hashPassword(salt, password)
	if err != nil {
		return "", err
	}

	return combine(salt, string(hash)), nil
}

func getPWPieces(password string) (salt string, hash string) {
	salt = password[:SaltLength]
	hash = password[SaltLength:len(password)]
	return
}

// IsMatch checks whether or not the correct password has been provided
func IsMatch(guess string, hash string) (bool, error) {
	salt, hash := getPWPieces(hash)
	saltedGuess := combine(salt, guess)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedGuess))
	return err == nil, err
}
