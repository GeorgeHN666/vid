package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	Conn     *mongo.Client
	Database string
}

func InitDatabase(Database, URI string) *db {

	ctx := context.TODO()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil
	}

	fmt.Println("<----- Successfully connected with database ----->")

	return &db{
		Conn:     c,
		Database: Database,
	}
}
