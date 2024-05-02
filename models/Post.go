package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"descrition" bson:"description"`
	LibraryK    string             `json:"libraryk" bson:"libraryk"`
	APIK        string             `json:"APIK" bson:"APIK"`
	Content     []VideoModel       `json:"content" bson:"content"`
}
