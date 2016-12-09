package server

import (
  "fmt"

    "github.com/kataras/iris"
)

type User struct {
  Username string `json:"username"`
    Password string `json:"password"`
    Apikey int `json:"api_key"`
}
var users map[string]User

func CreateServer() {
  users = make(m)
  iris.Post("webchat/signup", func(ctx *iris.Context) {
    user := &User{}
    if err := ctx.ReadJSON(user); err != nil {
      panic(err.Error())
    } else {
      if m[user.Username] != nil {
        m[user.Username] = user
        ctx.Write("Registered: %#v", user.Username)
      } else {
        ctx.Write("Username Exists")
      }
    }
  })

  iris.Listen("localhost:5700")
}
