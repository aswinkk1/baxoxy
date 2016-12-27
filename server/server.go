package server

import (
	// Third party packages
	"github.com/aswinkk1/baxoxy/controllers"
	"github.com/aswinkk1/baxoxy/jwthandler"
	//"github.com/dgrijalva/jwt-go"
	//jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	//"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	//"fmt"
	"log"
	"github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
)

func CreateServer() {

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	user := "gordon"
    pass := "secret!"
	router := fasthttprouter.New()
    router.POST("/webchat/signup", uc.CreateUser)
    router.POST("/webchat/login", uc.Login)
    router.GET("/",uc.Chathandler)
    router.GET("/webchat/protected/", jwthandler.BasicAuth(uc.Protected, user, pass))

    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
