package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

type MongoConn struct {
	client *mongo.Client
}

func (obj *MongoConn) Init(connectionUri string) {
	if obj.client != nil {
		return
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionUri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	obj.client = client

	if err != nil {
		log.Panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panic(err)
	}
}

func (obj *MongoConn) GetCollection(database string, name string) *mongo.Collection {
	return obj.client.Database(database).Collection(name)
}
