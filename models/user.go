package models

import "gopkg.in/mgo.v2/bson"

type (
	// User represents the structure of our resource
	User struct {
		Id     bson.ObjectId `json:"id" bson:"_id"`
		Username string `json:"username" bson:"username"`
		Password string `json:"password" bson:"password"`
		Apikey string `json:"api_key" bson:"api_key"`
	}
)