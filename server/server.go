package server

import (
	// Third party packages
	"github.com/aswinkk1/baxoxy/controllers"
	//"github.com/aswinkk1/baxoxy/jwthandler"
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

	/*myJwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	// Create a new user
	iris.Post("webchat/signup", uc.CreateUser)

	//login
	iris.Post("webchat/login", uc.Login)

	//test
	iris.Post("webchat/signin", myJwtMiddleware.Serve, uc.SecuredPingHandler)

	//	// Remove an existing user
	//	iris.DELETE("webchat/users/:id", uc.RemoveUser)
	iris.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.Write("CUSTOM 500 INTERNAL SERVER ERROR PAGE")
		// or ctx.Render, ctx.HTML any render method you want
		ctx.Log("http status: 500 happened!")
	})

	iris.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Write("CUSTOM 404 NOT FOUND ERROR PAGE")
		ctx.Log("http status: 404 happened!")
	})

	// emit the errors to test them
	iris.Get("/500", func(ctx *iris.Context) {
		ctx.EmitError(iris.StatusInternalServerError) // ctx.Panic()
	})

	iris.Get("/404", func(ctx *iris.Context) {
		ctx.EmitError(iris.StatusNotFound) // ctx.NotFound()
	})

	iris.Config.Websocket.Endpoint = "/my_endpoint"

	var myChatRoom = "room1"
    iris.Websocket.OnConnection(func(c iris.WebsocketConnection) {
    	c.Join(myChatRoom)

        c.On("chat", func(message string) {
        	c.To(myChatRoom).Emit("chat", "From: "+c.ID()+": "+message)
        })
        c.OnDisconnect(func() {
            fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
        })
    })

	// Fire up the server
	iris.Listen("localhost:5700")*/
	router := fasthttprouter.New()
    router.POST("/webchat/signup", uc.CreateUser)
    router.POST("/webchat/login", uc.Login)

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
