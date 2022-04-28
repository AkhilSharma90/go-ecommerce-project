package database

//
import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBSet is more Go-like compared to DBSet. Take a look at https://pkg.go.dev/database/sql@go1.17.2#OpenDB for example. Abbreviations are written in allcaps in Go. Examples: http.ServeTLS(), time.UTC
func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}

	// This doesn't work, even if you don't have mongodb running this line still prints.
	fmt.Println("Successfully Connected to the mongodb")
	return client
}

// This is a global var, don't use global vars for database connections.
// Instead use dependency injection to give your HTTP handlers access to
// the mongodb connection pool.
var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return collection

}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return productcollection
}
