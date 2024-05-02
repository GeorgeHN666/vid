package db

import (
	"context"
	"time"
	"video-streaming/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *db) InsertUser(user models.User) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("User")

	id := primitive.NewObjectID()
	user.ID = id

	_, err := db.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}

func (d *db) GetUser(email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("User")

	filter := bson.M{
		"email": bson.M{"$eq": email},
	}

	var res models.User

	err := db.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
