package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	pb "serverGRPC/server"
)

type server struct {
	pb.UnimplementedGetInfoServer
	rdb *redis.Client
}

type Data struct {
	Texto string
	Pais  string
}

var ctx = context.Background()

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	tweet := Data{
		Texto: in.GetTexto(),
		Pais:  in.GetPais(),
	}

	// Incrementar contador en Redis
	err := s.rdb.HIncrBy(ctx, "countries", tweet.Pais, 1).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println(tweet)

	return &pb.ReplyInfo{Info: "Hola cliente, recibí el album"}, nil
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

	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{rdb: rdb})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
