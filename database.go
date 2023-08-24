package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Database
var dbCtx = context.TODO()

func initDB() {

	clientOptions := options.Client().ApplyURI(loadedConfig.MongoDBUri)
	client, err := mongo.Connect(dbCtx, clientOptions)
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping(dbCtx, nil)
	if err != nil {
		log.Panic(err)
	}

	dbClient = client.Database(loadedConfig.MongoDBName)
}

func findOne(ctx context.Context, query interface{}, collectionName string) (bson.M, error) {
	collection := dbClient.Collection(collectionName)
	var result bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, echo.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func insertOne(ctx context.Context, collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := dbClient.Collection(collectionName)
	result, err := collection.InsertOne(ctx, document)
	return result, err
}

func upsertOne(ctx context.Context, collectionName string, filter interface{}, update interface{}) error {
	opts := options.Update().SetUpsert(true)
	collection := dbClient.Collection(collectionName)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if result.MatchedCount != 0 {
		fmt.Println("\nmatched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("\ninserted a new document with ID %v\n", result.UpsertedID)
		return nil
	}
	return nil
}

func updateOne(ctx context.Context, collectionName string, filter interface{}, update interface{}) error {
	collection := dbClient.Collection(collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
