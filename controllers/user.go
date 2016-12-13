package controllers

import (
	// "encoding/json"
	"log"
	"time"

	"github.com/aswinkk1/baxoxy/models"
	"github.com/aswinkk1/baxoxy/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}

	Response struct {
		Status  int    `json:"status"`
		Action  string `json:"action"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(ctx *iris.Context) {

	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	log.Println("Quer--", response)
	if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("user.Username", user.Username)
		if count, err := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).Count(); count == 0 {
			user.Id = bson.NewObjectId()
			uc.session.DB("baxoxy").C("users").Insert(user)
			response.Status = 200
			response.Action = "signup"
			response.Message = "Sign Up Successfull"
		} else {
			log.Println("Query--", response, " ", err)
			response.Status = 201
			response.Message = "Username already exists"
		}
	}
	log.Println("Quer--", response)
	ctx.JSON(iris.StatusCreated, response)
}

// Login removes an existing user resource
func (uc UserController) Login(ctx *iris.Context) {
	// Stub an user to be populated from the body
	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("user.Username", user.Username)
		dbData := models.User{}
		if error := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).One(&dbData); error != nil {
			log.Println(error.Error())
		} else {
			log.Println("db", dbData.Password, "APi", user.Password)
			if dbData.Password == user.Password {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"foo": "bar",
					"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
				})

				// Sign and get the complete encoded token as a string using the secret
				tokenString, err := token.SignedString([]byte("secret"))
				response.Status = 201
				response.Action = "login"
				response.Message = "user signed"
				response.Token = tokenString
				
				dbClient := services.RedisClient()
				dbErr := dbClient.Set(user.Username, tokenString, 0).Err()
				if dbErr != nil {
					log.Println("error", dbErr)
					response.Message ="Error"
				}
				log.Println("tokenString", tokenString, err)
			}
		}
	}
	ctx.JSON(iris.StatusCreated, response)
}


func (uc UserController) Logout(ctx *iris.Context) {
	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	if err := ctx.ReadJSON(&user); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("user.Username", user.Username)
		if count, err := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).Count(); count == 0 {
			log.Println(err.Error())
			response.Message ="Logged out Error"
		} else {			
			dbClient := services.RedisClient()
			err := dbClient.Del(user.Username).Err()
			if err != nil {
				response.Message ="Logged out Error"
			} else {
				response.Status = 200
				response.Action = "logout"
				response.Message = "logged Out successfully"
			}
			log.Println("token deleted with key", user.Username)
		}
	}
	ctx.JSON(iris.StatusCreated, response)
}

func (uc UserController) SecuredPingHandler(ctx *iris.Context) {
	
	ctx.Write("All good. You only get this message if you're authenticated")
}
