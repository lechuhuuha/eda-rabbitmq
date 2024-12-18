package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/lechuhuuha/eda-rabbitmq/internal"
)

var (
	exchangeEvent = "customer_events"
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

	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	if err := client.CreateExchange(exchangeEvent, "topic", false, false); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_created", "customers.created.*", exchangeEvent); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_test", "customers.*", exchangeEvent); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Send(ctx, exchangeEvent, "customers.created.us", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(`Test message between services`),
	}); err != nil {
		panic(err)
	}

	if err := client.Send(ctx, exchangeEvent, "customers.created.test", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Transient,
		Body:         []byte(`Test message between services undurable`),
	}); err != nil {
		panic(err)
	}

	defer client.Close()

	time.Sleep(20 * time.Second)
	fmt.Println(client)
}
