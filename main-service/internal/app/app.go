package app

import (
	"github.com/EvgeniyBudaev/go-base/main-service/internal/config"
	"github.com/EvgeniyBudaev/go-base/main-service/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type App struct {
	Logger logger.Logger
	config *config.Config
	db     *Database
	fiber  *fiber.App
}

func NewApp() *App {
	// Default logger
	defaultLogger, err := logger.NewLogger(logger.GetDefaultLevel())
	if err != nil {
		log.Fatal("error func NewApp, method NewLogger by path internal/app/app.go", err)
	}
	// Config
	cfg, err := config.Load(defaultLogger)
	if err != nil {
		log.Fatal("error func NewApp, method Load by path internal/app/app.go", err)
	}
	// Logger level
	loggerLevel, err := logger.NewLogger(cfg.LoggerLevel)
	if err != nil {
		log.Fatal("error func NewApp, method NewLogger by path internal/app/app.go", err)
	}
	// Database connection
	postgresConnection, err := newPostgresConnection(cfg)
	if err != nil {
		log.Fatal("error func NewApp, method newPostgresConnection by path internal/app/app.go", err)
	}
	database := NewDatabase(loggerLevel, postgresConnection)
	err = postgresConnection.Ping()
	if err != nil {
		log.Fatal("error func NewApp, method NewDatabase by path internal/app/app.go", err)
	}
	// Fiber
	f := fiber.New(fiber.Config{
		ReadBufferSize: 16384,
	})
	// CORS
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, X-Requested-With, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// GRPC
	// определяем порт для сервера
	//listen, err := net.Listen("tcp", ":3200")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// создаём gRPC-сервер без зарегистрированной службы
	//s := grpc.NewServer()
	//// регистрируем сервис
	//pb.RegisterUsersServer(s, &UsersServer{})
	//fmt.Println("Сервер gRPC начал работу")
	//// получаем запрос gRPC
	//if err := s.Serve(listen); err != nil {
	//	log.Fatal(err)
	//}

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

	log.Printf("Starting server service on port %s\n", cfg.Port)
	return &App{
		config: cfg,
		db:     database,
		Logger: loggerLevel,
		fiber:  f,
	}
}
