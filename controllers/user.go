package controllers

import (
	"encoding/json"
	"log"
	"time"
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

	Werror struct {
		Type  string `json:"type"`
		Data string `json:"data"`
	}
)

var Uc = NewUserController(getSession())

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
					response.Token = token
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

var tokenString string

type Message struct {
        To    string `json:"to"`
        Msg string `json:"msg"`
}

type Reply struct {
	Type string `json:"type"`
	Data Datas `json:"data"`
}
type Datas struct {
	Time string `json:"time"`
	Text string `json:"text"`
	To string `json:"to"`
	Author string `json:"author"`
}

func (uc UserController) Chathandler(ctx *fasthttp.RequestCtx) {
	tokenString = string(ctx.FormValue("token"))
    fmt.Println("Websocket request!\n",  tokenString)
    username, error :=jwthandler.TokenParser(tokenString)
    if value,ok := ActiveClients[username]; ok{
    	log.Println("close cheyyunna sthalam",value.websocket)
    	value := ActiveClients[username].websocket
    	value.Close();
    	deleteClient(username)
    }
    log.Println("username:", username, error)
    err := upgrader.Upgrade(ctx)
    log.Println(err)
}

var upgrader = websocket.Custom(Uc.chat,1,1)

func(uc UserController) chat(ws *websocket.Conn) {
	ws.SetReadLimit(2000)
    client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client, tokenString}
	addClient(sockCli)
	log.Println("ActiveClients",ActiveClients)
	for {
			//log.Println(len(ActiveClients), ActiveClients)
			var msg Message

			err := ws.ReadJSON(&msg)
			log.Println(msg)
			username,_ := RemoteUserMap[ws.RemoteAddr()]
			if err != nil {
				log.Println("msg.To:",msg.To)
				deleteClient(msg.To)
				log.Println(ActiveClients)
				log.Println("bye")
				log.Println(err)
				rep := Werror{Type: "alert", Data: "The message length exceeded" }
				if err := ActiveClients[username].websocket.WriteJSON(rep); err != nil{
					return
				}
				return
			}
			log.Println("RemoteAddrNow",ws.RemoteAddr())
			uc.broadcastMessage(msg, username)
	}
}

var ActiveClients = make(map[string]ClientConn)
var RemoteUserMap = make(map[net.Addr]string)
var ActiveClientsRWMutex sync.RWMutex
var RemoteUserMapRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
	token string
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	RemoteUserMapRWMutex.Lock()
	username, error :=jwthandler.TokenParser(tokenString)
	if error == nil{
		ActiveClients[username] = cc
		RemoteUserMap[cc.clientIP] = username
	}
	ActiveClientsRWMutex.Unlock()
	RemoteUserMapRWMutex.Unlock()
}

func deleteClient(cc string) {
	ActiveClientsRWMutex.Lock()
	//RemoteUserMapRWMutex.Lock()
	log.Println("ccdzs",cc)
	delete(ActiveClients, cc)
	//delete(RemoteUserMap, cc.clientIP)
	ActiveClientsRWMutex.Unlock()
	//RemoteUserMapRWMutex.Unlock()
}

func (uc UserController) broadcastMessage(msg Message,from string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()
	dbData := models.User{}
	if _,ok := ActiveClients[msg.To]; !ok{
		if error := uc.session.DB("baxoxy").C("users").Find(bson.M{"username": msg.To}).One(&dbData); error != nil {
			log.Println(error.Error(),"user not in db, user not exist")
			rep := Werror{Type: "alert", Data: "The user doesn't exists" }
			if err := ActiveClients[from].websocket.WriteJSON(rep); err != nil{
				return
			}
		}else{
			log.Println("userexist in db, user in offline")
			rep := Werror{Type: "alert", Data: "The user is not available" }
			if err := ActiveClients[from].websocket.WriteJSON(rep); err != nil{
				return
			}
		}
	}else{
		t := time.Now()
		var dat = Datas{Time:t.Format("2006/01/02/15:04:05"),Text:msg.Msg,To:msg.To,Author:from}
		var rep = Reply{Type:"message",Data: dat}
		if err := ActiveClients[msg.To].websocket.WriteJSON(rep); err != nil{
			return
		}
	}
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
