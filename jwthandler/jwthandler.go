package jwthandler

import (
    "fmt"
    "time"

    "github.com/dgrijalva/jwt-go"
)

const (
    mySigningKey = "WOW,MuchShibe,ToDogge"
)

func CreateToken(userName string) (string, error) {
    // Create the token
    mySigningKey := []byte("WOW,MuchShibe,ToDogge")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": userName,
        "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
    })
    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString(mySigningKey)
    return tokenString, err
}

func ParseToken(myToken string, myKey string) {
    token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
        return []byte(myKey), nil
    })

    if err == nil && token.Valid {
        fmt.Println("Your token is valid.  I like your style.")
    } else {
        fmt.Println("This token is terrible!  I cannot accept this.")
    }
}