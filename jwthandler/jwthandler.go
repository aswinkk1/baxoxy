package jwthandler

import (
    "fmt"
    "time"
    "log"
    "github.com/dgrijalva/jwt-go"
    "github.com/valyala/fasthttp"
    "encoding/json"
)

const (
    mySigningKey = "WOW,MuchShibe,ToDogge"
)
type MyCustomClaims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
type Response struct {
        Status  int    `json:"status"`
        Action  string `json:"action"`
        Message string `json:"message"`
}


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

func BasicAuth(h fasthttp.RequestHandler, requiredUser, requiredPassword string) fasthttp.RequestHandler {
    return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
        // Get the Basic Authentication credentials
        var response Response
        fmt.Println("basicauth")
        var myToken string
        if ctx.Request.Header.Peek("Authorization") != nil {
            myToken = string(ctx.Request.Header.Peek("Authorization"))
        }else{
            myToken = string(ctx.FormValue("token"))
        }
        myKey := "WOW,MuchShibe,ToDogge"
        log.Println(myKey)
        token, err := jwt.ParseWithClaims(myToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(myKey), nil
        })

        if err != nil {
            ctx.SetContentType("application/json")
            ctx.SetStatusCode(400)
            response.Status = 400
            response.Action = "websocket connection"
            response.Message = "The user doesn't exists"
            if b, err := json.Marshal(response); err == nil{
                fmt.Fprintf(ctx, string(b))
            }
        } else {
            if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
                fmt.Printf("%v", claims.Username)
                h(ctx)
            } else {
                fmt.Println(err)
            }
        }
    })
}

func TokenParser(myToken string) (string,error){
    myKey := "WOW,MuchShibe,ToDogge"
    log.Println("myKey:",myKey)
    var username string
    token, err := jwt.ParseWithClaims(myToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(myKey), nil
    })

    if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
        fmt.Printf("%v", claims.Username)
        username = claims.Username
    } else {
        fmt.Println(err)
    }
    return username,err
}
