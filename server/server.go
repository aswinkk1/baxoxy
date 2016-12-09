package server

import (

    "github.com/kataras/iris"
)

type User struct {
  Username string `json:"username"`
    Password string `json:"password"`
    Apikey int `json:"api_key"`
}
var users map[string]User

func CreateServer() {
  users = make(map[string]User)
  iris.Post("webchat/signup", func(ctx *iris.Context) {
    user := &User{}
    if err := ctx.ReadJSON(user); err != nil {
      panic(err.Error())
    } else {
      ctx.Write("Registered: %v", user.Username)
      if usr, ok := users[user.Username]; ok == false {
        users[user.Username] = *user
        ctx.Write("Registered: %#v", user.Username)
      } else {
        ctx.Write("Username Exists: %v", usr.Username)
      }
    }
  })
  iris.Get("webchat/users", func(ctx *iris.Context) {
    ctx.JSON(iris.StatusOK, users)
  })

  iris.Listen("localhost:5700")
}
