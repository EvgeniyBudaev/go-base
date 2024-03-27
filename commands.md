Инициализация зависимостей

```
go mod init github.com/EvgeniyBudaev/go-base/main-service
go mod init github.com/EvgeniyBudaev/go-base/grpc-service
go mod init github.com/EvgeniyBudaev/go-base/rabbitmq-service
```

Библиотека для работы с переменными окружения ENV
https://github.com/joho/godotenv

```
go get -u github.com/joho/godotenv
```

ENV Config
https://github.com/kelseyhightower/envconfig

```
go get -u github.com/kelseyhightower/envconfig
```

Логирование
https://pkg.go.dev/go.uber.org/zap

```
go get -u go.uber.org/zap
```

Fiber
https://github.com/gofiber/fiber

```
go get -u github.com/gofiber/fiber/v2
```

CORS
https://github.com/gorilla/handlers

```
go get -u github.com/gorilla/handlers
```

Подключение к БД
Драйвер для Postgres

```
go get -u github.com/lib/pq
```

Вызовите утилиту protoc для генерации соответствующих go-файлов. Для этого выполните команду:
```
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/demo.proto
```

В --go-out запишется файл с кодом для Protobuf-сериализации. В --go-grpc_out сохранится файл с gRPC-интерфейсами и методами. Так как вы указали параметр paths=source_relative, сгенерированные файлы создадутся в поддиректории ./proto. Если бы указали параметр paths=import, то сгенерированные файлы создались бы в директории, указанной в директиве go_package, то есть ./demo/proto.

gRPC Protocol Buffer Compiler Installation https://grpc.io/docs/protoc-installation/
```
sudo apt install -y protobuf-compiler
protoc --version
```

После этого установите утилиты, которые отвечают за кодогенерацию go-файлов:
```
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
```

Разработка gRPC-сервера После того как были реализованы все необходимые интерфейсы,
можно приступать к созданию функции main. Она запустит gRPC-сервер.

Вот алгоритм по шагам:

При вызове net.Listen указать порт, который будет прослушивать сервер.
Создать экземпляр gRPC-сервера функцией grpc.NewServer().
Зарегистрировать созданный сервис UsersServer на сервере gRPC.
Вызвать Serve() для начала работы сервера. Он будет слушать указанный порт, пока процесс не прекратит работу.
Разработка gRPC-клиента Соединение с сервером устанавливается при вызове функции grpc.Dial().
В первом параметре указывается адрес сервера, далее перечисляются опциональные параметры.
Функция pb.NewUsersClient(conn) возвращает переменную интерфейсного типа UsersClient,
для которого сгенерированы методы с соответствующими запросами из proto-файла.


Запуск RabbitMQ
```
docker-compose up
```
15672 - порт админки
5672 - порт RabbitMQ
После запуска можем перейти в админку по адресу http://localhost:15672

https://github.com/rabbitmq/rabbitmq-tutorials/tree/main/go
```
go get github.com/rabbitmq/amqp091-go
```

Stop process
```
sudo lsof -i :15672
sudo lsof -i :5432
sudo kill PID_number
```
