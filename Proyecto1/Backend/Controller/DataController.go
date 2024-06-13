package Controller

import (
	"Backend/Instance"
	"Backend/Model"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// Insertar

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

// Eliminar

func DeleteDataRAM(nameCol string) {
	collection := Instance.Mg.Db.Collection(nameCol)
	filter := bson.D{}

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if count >= 20 {
		_, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Documentos eliminados de coleccion", nameCol)
	}
}

func DeleteDataCPU(nameCol string) {
	collection := Instance.Mg.Db.Collection(nameCol)
	filter := bson.D{}

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if count >= 20 {
		_, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Documentos eliminados de coleccion", nameCol)
	}
}

func DeleteDataProcess(nameCol string) {
	collection := Instance.Mg.Db.Collection(nameCol)
	filter := bson.D{}

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if count >= 14880 {
		_, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Documentos eliminados de coleccion", nameCol)
	}
}
