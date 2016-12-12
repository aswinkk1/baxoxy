package controllers

import (
	"encoding/json"
	"log"
	"time"
	
	"github.com/dgrijalva/jwt-go"
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
	
	Response struct {
		status int
		action string
		message	string
		token string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(ctx *iris.Context) {
	
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

// Login removes an existing user resource
func (uc UserController) Login(ctx *iris.Context) {
	// Stub an user to be populated from the body
	user := models.User{}
	response := `{"status" : 400, "error_message" :"Error"}`
	
    if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
	} else { 
		log.Println("user.Username", user.Username)
		dbData := models.User{}
		if error := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).One(&dbData); error != nil{
			log.Println(error.Error())
		} else {
			log.Println("db",dbData.Password, "APi", user.Password)
			if dbData.Password == user.Password {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"foo": "bar",
				"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
				tokenString, err := token.SignedString([]byte("secret"))
				response = `{ "status": 201, "action": "login", "message": "user signed", "token":"" }`	
				log.Println("tokenString", tokenString, err)
			}
		}
	}
	ctx.JSON(iris.StatusCreated, response)
}

func (uc UserController) SecuredPingHandler(ctx *iris.Context) {
    ctx.Write("All good. You only get this message if you're authenticated")
}