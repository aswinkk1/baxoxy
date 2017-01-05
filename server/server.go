package server

import (
	// Third party packages
	"github.com/aswinkk1/baxoxy/controllers"
	"github.com/aswinkk1/baxoxy/jwthandler"
	//"github.com/dgrijalva/jwt-go"
	//jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	//"github.com/kataras/iris"
	//"gopkg.in/mgo.v2"
	//"fmt"
	"log"
	"github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
)

func CreateServer() {

	// Get a UserController instance
	uc := controllers.Uc
	user := "gordon"
    pass := "secret!"
	router := fasthttprouter.New()
    router.POST("/webchat/signup", uc.CreateUser)
    router.POST("/webchat/login", uc.Login)
    router.GET("/", jwthandler.BasicAuth(uc.Chathandler, user, pass))
    router.GET("/webchat/protected/", jwthandler.BasicAuth(uc.Protected, user, pass))

    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

