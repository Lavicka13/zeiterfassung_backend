package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("geheim") // später aus .env holen!

type Claims struct {
    NutzerID  uint
    RechteID  uint
    Email     string
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
    return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}
