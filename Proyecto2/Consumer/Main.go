package main

import (
    "Consumer/Database"
    "Consumer/Instance"
    "Consumer/model"
    "context"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    "github.com/segmentio/kafka-go"
    "log"
    "time"
)

func processEvent(event []byte) {
    var data model.Data
    err := json.Unmarshal(event, &data)
    if err != nil {
        log.Fatal(err)
    }

    // Conectar a la base de datos si no está conectado
    if Instance.Mg.Client == nil {
        if err := Database.Connect(); err != nil {
            log.Fatal("Error en", err)
        }
    }

    // Establecer los campos de fecha y hora
    data.CreatedAt = time.Now()

    collection := Instance.Mg.Db.Collection("register")
    _, err = collection.InsertOne(context.TODO(), data)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    topic := "mytopic"
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     []string{"my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"},
        Topic:       topic,
        Partition:   0,
        MinBytes:    10e3,
        MaxBytes:    10e6,
        StartOffset: kafka.LastOffset,
        GroupID:     uuid.New().String(),
    })

    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message: %v", err)
        }
        fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

        processEvent(m.Value)

        err = r.CommitMessages(context.Background(), m)
        if err != nil {
            log.Printf("Error committing message: %v", err)
        }
    }
}
