package controllers

import (
//	"encoding/json"
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
	response := `{"status" : 400, "error_message" :"Error"}`
	
    if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
    } else {	
		log.Println("user.Username", user.Username)
		if count, err := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).Count(); count == 0 {
		user.Id = bson.NewObjectId()
		uc.session.DB("baxoxy").C("users").Insert(user)
		response = `{ "status": 200, "action": "signup", "message": "Sign Up Successful" }`
		} else {
			log.Println("Query--", count," ", err)
			response = `{ "status": 201, "action": "signup", "message": "username already exists" }`	
		}
    }
	ctx.JSON(iris.StatusCreated, response)
}

// RemoveUser removes an existing user resource