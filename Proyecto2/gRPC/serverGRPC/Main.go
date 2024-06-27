package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"serverGRPC/kafka"
	"serverGRPC/model"
	pb "serverGRPC/server"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGetInfoServer
	rdb           *redis.Client
	totalMessages int64 // Contador global para el total de mensajes
	mu            sync.Mutex
}

type Data struct {
	Texto string
	Pais  string
}

var ctx = context.Background()

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	tweet := model.Data{
		Texto: in.GetTexto(),
		Pais:  in.GetPais(),
	}

	// Procesar la solicitud recibida
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalMessages++

	err := s.rdb.HIncrBy(ctx, "countries", in.GetPais(), 1).Err()
	if err != nil {
		return nil, err
	}

	// Incrementar el contador global y enviar a Redis con clave "total_messages"
	err = s.rdb.Set(ctx, "total_messages", s.totalMessages, 0).Err()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Texto: %s, País: %s\n", in.GetTexto(), in.GetPais())

	kafka.Produce(tweet)

	// Devolver la respuesta con los datos procesados
	return &pb.ReplyInfo{
		Texto: in.GetTexto(),
		Pais:  in.GetPais(),
	}, nil
}

func main() {
	// Inicializar el cliente Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Probar la conexión con Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("No se pudo conectar a Redis: %v", err)
	}

	// Configurar el servidor gRPC para escuchar en el puerto 3001
	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		panic(err)
	}

	// Crear una instancia del servidor gRPC
	s := grpc.NewServer()

	// Registrar el servicio gRPC generado y el servidor personalizado
	pb.RegisterGetInfoServer(s, &server{rdb: rdb})

	// Iniciar el servidor gRPC
	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
