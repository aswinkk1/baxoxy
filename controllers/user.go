package controllers

import (
	"encoding/json"
	"log"
	//"time"
	"bytes"
	"github.com/aswinkk1/baxoxy/models"
	"github.com/aswinkk1/baxoxy/jwthandler"
	"github.com/aswinkk1/baxoxy/password"
	//"github.com/dgrijalva/jwt-go"
	//"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
    "fmt"
    "github.com/fasthttp-contrib/websocket"
    "net"
    "sync"

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
func (uc UserController) CreateUser(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	s :=   ctx.PostBody()
	postbody := bytes.NewBuffer(s)
	log.Println("postbody\n", postbody)
	err := json.NewDecoder(postbody).Decode(&user)
	if err != nil {
	    log.Println("error:", err)
	}else{
		pass := libs.Password{}
		user.Password = pass.Gen(string(user.Password))
		if count, err := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).Count(); count == 0 {
			user.Id = bson.NewObjectId()
			uc.session.DB("baxoxy").C("users").Insert(user)
			log.Println("usercreated")
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(200)
			response.Status = 200
			response.Action = "signup"
			response.Message = "user created"
			if b, err := json.Marshal(response); err == nil{
				fmt.Fprintf(ctx, string(b))
			}
		} else {
			log.Println("userAlreadyexist",err)
			ctx.SetContentType("application/json")
			response.Status = 400
			response.Action = "signup"
			response.Message = "user already exist"
			if b, err := json.Marshal(response); err == nil{
				ctx.Error(string(b), 400)
			}

		}
	}
}

// Login removes an existing user resource
func (uc UserController) Login(ctx *fasthttp.RequestCtx) {
	log.Println("Login")
	// Stub an user to be populated from the body
	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	s :=   ctx.PostBody()
	postbody := bytes.NewBuffer(s)
	log.Println("postbody\n", postbody)
	err := json.NewDecoder(postbody).Decode(&user)
	if err != nil {
	    log.Println("error:", err)
	    ctx.SetContentType("application/json")
			response.Status = 400
			response.Action = "login"
			response.Message = "log in failed"
			if b, err := json.Marshal(response); err == nil{
				ctx.Error(string(b), 400)
			}
	}else{
		dbData := models.User{}
		if error := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": user.Username}).One(&dbData); error != nil {
			log.Println(error.Error())
		} else {
			log.Println("db", dbData.Password, "APi", user.Password)
			pass := libs.Password{}
			var cp = pass.Compare(dbData.Password, user.Password)
			log.Println("resp",cp)
			if cp {
				if token, err := jwthandler.CreateToken(user.Username); err == nil{
					log.Println("token",token)
					ctx.SetContentType("application/json")
					ctx.SetStatusCode(200)
					response.Status = 200
					response.Action = "login"
					response.Message = "login successfull"
					if b, err := json.Marshal(response); err == nil{
						fmt.Fprintf(ctx, string(b))
					}
				}
			}
		}
	}
}

func (uc UserController) Protected(ctx *fasthttp.RequestCtx) {
    fmt.Println("Protected!\n")
}

func (uc UserController) Chathandler(ctx *fasthttp.RequestCtx) {
    fmt.Println("Websocket request!\n")
    err := upgrader.Upgrade(ctx)
    log.Println(err)
}

var upgrader = websocket.New(chat)

func chat(ws *websocket.Conn) {
    client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client}
	addClient(sockCli)

	for {
			log.Println(len(ActiveClients), ActiveClients)
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				deleteClient(sockCli)
				log.Println("bye")
				log.Println(err)
				return
			}
			broadcastMessage(messageType, p)
	}
}

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}