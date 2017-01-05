package server

import (
	"log"
	// Third party packages
	"github.com/aswinkk1/baxoxy/controllers"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
)
func WebSocket(){
	iris.Config.Websocket.Endpoint = "/test"
    // for Allow origin you can make use of the middleware
    //iris.Config.Websocket.Headers["Access-Control-Allow-Origin"] = "*"

    var myChatRoom = "room1"
    iris.Websocket.OnConnection(func(c iris.WebsocketConnection) {
		spew.Dump(c)
        c.Join(myChatRoom)
		log.Println("\nConnection with ID: %s has been connected!", c.ID())
        c.On("chat", func(message string) {
            c.To(myChatRoom).Emit("chat", "From: "+c.ID()+": "+message)
        })

        c.OnDisconnect(func() {
            log.Println("\nConnection with ID: %s has been disconnected!", c.ID())
        })
    })
}


func CreateServer() {

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	uc.SetupDb()

	myJwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	// Create a new user
	iris.Post("webchat/signup", uc.CreateUser)

	//login
	iris.Post("webchat/login", uc.Login)

	//logout
	iris.Post("webchat/logout", uc.Logout)


	//test
	iris.Post("webchat/signin", myJwtMiddleware.Serve, uc.SecuredPingHandler)

	iris.Config.Websocket.Endpoint = "/"
    // for Allow origin you can make use of the middleware
    //iris.Config.Websocket.Headers["Access-Control-Allow-Origin"] = "*"
	 var myChatRoom = "room1"
    iris.Websocket.OnConnection(func(c iris.WebsocketConnection) {
		//log.Println("params", c)
		//spew.Dump(c)
        c.Join(myChatRoom)
		log.Println("\nConnection with ID: %s has been connected!", c.ID())
        c.On("chat", func(message string) {
			log.Println("From: ", c.ID(), ":message ", message)
            c.To(myChatRoom).Emit("chat", "From: "+c.ID()+": "+message)
        })

        c.OnDisconnect(func() {
            log.Println("\nConnection with ID: %s has been disconnected!", c.ID())
        })
    })
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

	// Fire up the server
	iris.Listen("localhost:5700")
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
