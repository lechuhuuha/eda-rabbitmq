package main

import (
	"fmt"
	"time"

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

	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	if err := client.CreateExchange("customer_events", "topic", false, false); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_created", "customers.created.*", "customer_events"); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_test", "customers.*", "customer_events"); err != nil {
		panic(err)
	}

	defer client.Close()

	time.Sleep(20 * time.Second)
	fmt.Println(client)
}
