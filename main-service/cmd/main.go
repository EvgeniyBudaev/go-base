package main

import (
	"context"
	"fmt"
	"github.com/EvgeniyBudaev/go-base/main-service/internal/app"
	pb "github.com/EvgeniyBudaev/go-base/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer cancelCtx()

	var wg sync.WaitGroup
	application := app.NewApp()

	wg.Add(1)
	go func() {
		// GRPC
		// определяем порт для сервера
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}
		// создаём gRPC-сервер без зарегистрированной службы
		s := grpc.NewServer()
		// регистрируем сервис
		pb.RegisterUsersServer(s, &app.UsersServer{})
		fmt.Println("Сервер gRPC начал работу")
		// получаем запрос gRPC
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// RabbitMQ
		// Consumer
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
		}
		defer func() {
			_ = conn.Close() // Закрываем подключение в случае удачной попытки подключения
		}()
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("failed to open a channel. Error: %s", err)
		}
		defer func() {
			_ = ch.Close() // Закрываем подключение в случае удачной попытки подключения
		}()
		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			log.Fatalf("failed to declare a queue. Error: %s", err)
		}
		messages, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}
		var forever chan struct{}
		go func() {
			for message := range messages {
				log.Printf("received a message: %s", message.Body)
			}
		}()
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever

		log.Printf("Starting server service on port %s\n", application.Config.Port)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// REST API
		if err := application.StartHTTPServer(ctx); err != nil {
			application.Logger.Fatal("error func main, method StartHTTPServer by path cmd/main.go", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
}
