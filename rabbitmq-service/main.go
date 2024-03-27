// Producer
package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}
	defer func() {
		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
	}()
	queueName := "hello"
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("failed to publish a message. Error: %s", err)
	}
	log.Printf(" [x] Sent %s\n", body)

	//conn, err := amqp.Dial("amqp://guest:guest@:15672/")
	//if err != nil {
	//	log.Fatal("error func NewApp, method amqp.Dial by path internal/app/app.go", err)
	//}
	//defer conn.Close()
	//ch, err := conn.Channel()
	//if err != nil {
	//	log.Fatal("error func NewApp, method conn.Channel by path internal/app/app.go", err)
	//}
	//defer ch.Close()
	//queueName := "hello"
	//q, err := ch.QueueDeclare(
	//	queueName, // name
	//	false,     // durable
	//	false,     // delete when unused
	//	false,     // exclusive
	//	false,     // no-wait
	//	nil,       // arguments
	//)
	//if err != nil {
	//	log.Fatal("error func NewApp, method ch.QueueDeclare by path internal/app/app.go", err)
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//body := "Hello World!"
	//err = ch.PublishWithContext(ctx,
	//	"",     // exchange
	//	q.Name, // routing key
	//	false,  // mandatory
	//	false,  // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	//if err != nil {
	//	log.Fatal("error func NewApp, method ch.PublishWithContext by path internal/app/app.go", err)
	//}
	//log.Printf(" [x] Sent %s\n", body)
}
