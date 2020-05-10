package middleware

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ldXiao/GoReactChatApp/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const connectionString = "Connection String"

// collection object/instance
var UsersCollection *mongo.Collection

// collections object/instance
var ChatsCollection *mongo.Collection

// create connection with mongo db
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(config.ConnectionString)

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

	fmt.Println("Connected to MongoDB!")

	UsersCollection = client.Database(config.DbName).Collection(config.CollNameUsers)

	indexName, err := UsersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Fatal(err)
		fmt.Fprint(os.Stderr, "%q failed, duplicate email insied the collection\n", indexName)
	}

	ChatsCollection = client.Database(config.DbName).Collection(config.CollNameChats)

	fmt.Println("Collection instance created!")
}
