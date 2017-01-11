package server

import (
	"log"
	"strings"
	"encoding/json"
	"time"
	// Third party packages
	"github.com/aswinkk1/baxoxy/controllers"
	//"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
)

var ActiveClients = make(map[string]string)

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
	 //var myChatRoom = "room1"
    iris.Websocket.OnConnection(func(c iris.WebsocketConnection) {
    	id := c.ID()
    	log.Println("\nConnection with ID: %s has been connected!", id)
		token := strings.TrimPrefix(c.Request().RequestURI,"/?token=")
		if username,err := TokenParser(token); err == nil{
			ActiveClients[username] = id
			log.Println(ActiveClients)
		}else{
			go c.Disconnect();
			//log.Println("er",er)
			log.Println("err",err)
		}
		//spew.Dump(c)
        //c.Join(myChatRoom)
        c.OnMessage(func(message []byte){
        	var msg Message
        	log.Println(string(message))
        	if err := json.Unmarshal(message, &msg); err != nil {
        		panic(err)
    		}
    		token := strings.TrimPrefix(c.Request().RequestURI,"/?token=")
    		username,_ := TokenParser(token)
    		t := time.Now()
    		var dat = Datas{Time:t.Format("2006/01/02/15:04:05"),Text:msg.Msg,To:msg.To,Author:username}
			var rep = Reply{Type:"message",Data: dat}
			b, _ := json.Marshal(rep)
    		c.To(ActiveClients[msg.To]).EmitMessage(b)
    		c.EmitMessage(b)
        })

        c.OnDisconnect(func() {
        	token := strings.TrimPrefix(c.Request().RequestURI,"/?token=")
        	log.Println(token)
        	username,_ := TokenParser(token)
        	delete(ActiveClients,username)
            log.Println("\nConnection with ID: %s has been disconnected!", c.ID())
        })
    })
	//	// Remove an existing user
	//	iris.DELETE("webchat/users/:id", uc.RemoveUser)
	iris.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.Write([]byte("CUSTOM 500 INTERNAL SERVER ERROR PAGE"))
		// or ctx.Render, ctx.HTML any render method you want
		ctx.Log("http status: 500 happened!")
	})

	iris.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Write([]byte("CUSTOM 404 NOT FOUND ERROR PAGE"))
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
	iris.Listen("10.7.20.26:5700")
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

type MyCustomClaims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}


func TokenParser(myToken string) (string,error){
    myKey := "secret"
    var username string
    token, err := jwt.ParseWithClaims(myToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(myKey), nil
    })

    if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
        log.Printf("%v", claims.Username)
        username = claims.Username
    } else {
        log.Println(err)
    }
    return username,err
}
