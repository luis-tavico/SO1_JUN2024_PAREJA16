package Controller

import (
	"Backend/Instance"
	"Backend/Model"
	"context"
	"log"
)

func InsertData(nameCol string, data Model.Data) {
	collection := Instance.Mg.Db.Collection(nameCol)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
}
