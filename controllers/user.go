package controllers

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris"
	"github.com/aswinkk1/baxoxy/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}


// CreateUser creates a new user resource
func (uc UserController) CreateUser(ctx *iris.Context) {
	log.Println("calla")
	user := models.User{}
	
    if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
    } else {	
		// Populate the user data
//		json.NewDecoder(body).Decode(&user)
		
		// Add an Id
		user.Id = bson.NewObjectId()
		
		uc.session.DB("baxoxy").C("users").Insert(user)
		// Marshal provided interface into JSON structure
		uj, _ := json.Marshal(string(user))
        ctx.Write("Registered: %#v", uj)
      
    }	
}

// RemoveUser removes an existing user resource