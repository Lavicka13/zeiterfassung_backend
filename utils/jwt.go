package utils

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// getJWTKey ruft den JWT-Schlüssel aus der Umgebungsvariable ab oder verwendet einen Fallback-Wert
func getJWTKey() []byte {
    key := os.Getenv("JWT_SECRET")
    if key == "" {
        key = "geheim" // Fallback für Entwicklung, aber nicht für Produktion empfohlen!
    }
    return []byte(key)
}

type Claims struct {
    NutzerID  uint   `json:"nutzer_id"`
    RechteID  uint   `json:"rechte_id"`
    Email     string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(id uint, email string, rechte uint) (string, error) {
    claims := Claims{
        NutzerID: id,
        RechteID: rechte,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(getJWTKey())
}

func ParseToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return getJWTKey(), nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}