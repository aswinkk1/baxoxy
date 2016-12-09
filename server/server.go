package server

import (
	// Third party packages
	"github.com/kataras/iris"
	"github.com/aswinkk1/baxoxy/controllers"
	"gopkg.in/mgo.v2"
)

func CreateServer() {

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())
	
	//test
//	iris.Get("webchat/users", func(ctx *iris.Context) {
//		ctx.JSON(iris.StatusOK, {})
//  	})
	
//	// Get a user resource
//	iris.GET("webchat/users/:id", uc.GetUser)

	// Create a new user
	iris.Post("webchat/signup", uc.CreateUser)

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
