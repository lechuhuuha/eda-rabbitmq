package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/lechuhuuha/eda-rabbitmq/constant"
	"github.com/lechuhuuha/eda-rabbitmq/internal"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("lchh", "lchh-secret", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	consumeConn, err := internal.ConnectRabbitMQ("lchh", "lchh-secret", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer consumeConn.Close()

	consumeClient, err := internal.NewRabbitMQClient(consumeConn)
	if err != nil {
		panic(err)
	}
	defer consumeClient.Close()

	time.Sleep(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		if err := client.Send(ctx, constant.ExchangeEvent, "customers.created.us", amqp091.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp091.Persistent,
			Body:         []byte(`Test message between services`),
		}); err != nil {
			panic(err)
		}
	}

	defer client.Close()
	fmt.Println(client)
}
