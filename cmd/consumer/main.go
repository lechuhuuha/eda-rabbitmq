package main

import (
	"fmt"

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

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}
	var blockChan chan struct{}
	go func() {
		for message := range messageBus {
			fmt.Println("New message", message)
			if err := message.Ack(false); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Acked the message", message.AppId)
		}
	}()
	fmt.Println("Press CTRL + C to exit")

	<-blockChan
}
