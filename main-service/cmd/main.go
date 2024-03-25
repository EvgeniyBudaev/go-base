package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/EvgeniyBudaev/go-base/proto"
)

func main() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterUsersServer(s, &UsersServer{})
	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
