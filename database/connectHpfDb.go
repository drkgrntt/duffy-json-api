package database

import (
	"context"
	"fmt"
	"log"

	"github.com/drkgrntt/duffy-json-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	hpfDB *mongo.Database
)

func ConnectHpfDB(config *utils.Config) {
	var err error
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", config.HpfDBUser, config.HpfDBPassword, config.HpfDBHost, config.HpfDBName)

	fmt.Println("? Connecting using the following URI: " + uri)

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Failed to connect to the Mongo Database")
	}

	hpfDB = client.Database("dragonflyer-hotel-prices")
	fmt.Println("? Connected Successfully to the Mongo Database")
}

func GetHpfDatabase() *mongo.Database {
	return hpfDB
}
