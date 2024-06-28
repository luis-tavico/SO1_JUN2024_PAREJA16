package Database

import (
	"Consumer/Instance"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func Connect() error {
	server := "mongo-service"
	port := "27017"
	dbName := "DB2"
	var mongoUri = "mongodb://" + server + ":" + port

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		log.Fatal(err)
	}

	Instance.Mg = Instance.MongoInstance(MongoInstance{
		Client: client,
		Db:     db,
	})

	return nil
}
