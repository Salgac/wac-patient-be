package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName string
var collection string

func ConnectDB() *mongo.Client {
	// helper
	enviro := func(name string, defaultValue string) string {
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return defaultValue
	}

	// get values
	host := enviro("AMBULANCE_API_MONGODB_HOST", "localhost")
	port := enviro("AMBULANCE_API_MONGODB_PORT", "27017")
	if port, err := strconv.Atoi(port); err == nil {
	} else {
		log.Printf("Invalid port value: %v", port)
		port = 27017
	}
	userName := enviro("AMBULANCE_API_MONGODB_USERNAME", "")
	password := enviro("AMBULANCE_API_MONGODB_PASSWORD", "")
	dbName = enviro("AMBULANCE_API_MONGODB_DATABASE", "xsalgovic-patient-wl")
	collection = enviro("AMBULANCE_API_MONGODB_COLLECTION", "patients")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://" + userName + ":" + password + "@" + host + ":" + port)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

func GetCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(dbName).Collection(collection)
}
