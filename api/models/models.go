package models

import "time"

type Post struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Body      string    `json:"body" bson:"body"`
	UserName  string    `json:"user_name" bson:"user_name"`
}
