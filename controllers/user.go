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
    //"net"
    //"sync"

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
	// Stub an user to be populated from the body
	user := models.User{}
	response := Response{Status: 400, Message: "Error"}
	s :=   ctx.PostBody()
	postbody := bytes.NewBuffer(s)
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
			pass := libs.Password{}
			var cp = pass.Compare(dbData.Password, user.Password)
			if cp {
				if token, err := jwthandler.CreateToken(user.Username); err == nil{
					log.Println("token",token)
					ctx.SetContentType("application/json")
					ctx.SetStatusCode(200)
					response.Status = 200
					response.Action = "login"
					response.Message = "login successfull"
					log.Println("login successfull")
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
    /*if value,ok := ActiveClients[username]; ok{
    	log.Println("close cheyyunna sthalam",value.websocket)
    	value := ActiveClients[username].websocket
    	value.Close();
    	deleteClient(username)
    }*/
    log.Println("username:", username, error)
    err := upgrader.Upgrade(ctx)
    log.Println("websocket connected")
    log.Println(err)
}

var upgrader = websocket.Custom(Uc.chat,1,1)
var ActiveClients = make(map[string]*Client)

func(uc UserController) chat(ws *websocket.Conn) {
	client := &Client{
        ws:   ws,
        send: make(chan []byte),
    }
    username, _ :=jwthandler.TokenParser(tokenString)
    ActiveClients[username] = client
    hub.addClient <- client

    go client.write()
    client.read()
}


type Hub struct {
    clients map[*Client]bool
    broadcast     chan []byte
    addClient     chan *Client
    removeClient  chan *Client
}

// initialize a new hub
var hub = Hub{
    broadcast:     make(chan []byte),
    addClient:     make(chan *Client),
    removeClient:  make(chan *Client),
    clients:       make(map[*Client]bool),
}

// Runs forever as a goroutine
func (hub *Hub) start() {
    for {
        // one of these fires when a channel
        // receives data
        select {
        case conn := <-hub.addClient:
            // add a new client
            hub.clients[conn] = true
        case conn := <-hub.removeClient:
            // remove a client
            if _, ok := hub.clients[conn]; ok {
                delete(hub.clients, conn)
                close(conn.send)
            }
        case message := <-hub.broadcast:
            // broadcast a message to all clients
            var msg Message
            if err := json.Unmarshal(message, &msg); err != nil {
        		panic(err)
    		}
            for conn := range hub.clients {
                select {
                case conn.send <- message:
                	 cl := ActiveClients[msg.To]
                	 log.Println(cl)
                	 log.Println(hub.clients[cl])
                	 log.Println("conn",conn)
                default:
                    close(conn.send)
                    delete(hub.clients, conn)
                }
            }
        }
    }
}


type Client struct {
    ws *websocket.Conn
    // Hub passes broadcast messages to this channel
    send chan []byte
}

// Hub broadcasts a new message and this fires
func (c *Client) write() {
	var msg Message
    // make sure to close the connection incase the loop exits
    defer func() {
        c.ws.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.ws.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := json.Unmarshal(message, &msg); err != nil {
        		panic(err)
    		}
    		if(c == ActiveClients[msg.To]){
    			c.ws.WriteMessage(websocket.TextMessage, []byte(msg.To))
    		}

        }
    }
}

// New message received so pass it to the Hub
func (c *Client) read() {
    defer func() {
        hub.removeClient <- c
        c.ws.Close()
    }()

    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            hub.removeClient <- c
            c.ws.Close()
            break
        }

        hub.broadcast <- message
    }
}




// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {

	go hub.start()
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}

