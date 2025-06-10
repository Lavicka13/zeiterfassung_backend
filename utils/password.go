package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword erstellt ein bcrypt-Hash aus dem übergebenem Passwort
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

// CheckPasswordHash überprüft, ob das Passwort zum gegebenen Hash passt. Kombination gültig => true
func CheckPasswordHash(password, hash string) bool {
    // bcrypt vergleicht intern das gehashte Passwort mit dem Klartextpasswort.
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
