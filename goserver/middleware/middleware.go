package middleware

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection string
const connectionString = "mongodb://localhost:27017"

// const connectionString = "mongodb+srv://admin:1password1@cluster0-yauku.mongodb.net/test?retryWrites=true&w=majority"

// const connectionString = "Connection String"

// Database Name
const dbName = "test"

// Collection name
const collNameUsers = "Users1"

const collNameChats = "Chats1"

// collection object/instance
var UsersCollection *mongo.Collection

// collections object/instance
var ChatsCollection *mongo.Collection

// create connection with mongo db
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	UsersCollection = client.Database(dbName).Collection(collNameUsers)

	ChatsCollection = client.Database(dbName).Collection(collNameChats)

	fmt.Println("Collection instance created!")
}
