package db

import (
	"context"
	"time"
	"video-streaming/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *db) InsertNewPost(post *models.Post) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("post")

	_, err := db.InsertOne(ctx, post)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return post.ID, nil

}

func (d *db) GetPost(i string) (*models.Post, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("post")

	id, err := primitive.ObjectIDFromHex(i)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	var post models.Post

	err = db.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (d *db) GetArrayOfPosts() ([]*models.Post, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("post")

	var res []*models.Post

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var post models.Post

		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}

		res = append(res, &post)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *db) DeletePost(uid string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	db := d.Conn.Database(d.Database).Collection("post")

	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	_, err = db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
