package Controller

import (
	"Backend/Instance"
	"Backend/Model"
	"context"
	"log"
)

func InsertDataRAM(nameCol string, data Model.DataRAM) {
	collection := Instance.Mg.Db.Collection(nameCol)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertDataCPU(nameCol string, data Model.DataCPU) {
	collection := Instance.Mg.Db.Collection(nameCol)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertDataProcess(nameCol string, data Model.DataProcess) {
	collection := Instance.Mg.Db.Collection(nameCol)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
}
