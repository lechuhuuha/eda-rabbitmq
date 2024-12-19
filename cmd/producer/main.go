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
	conn, err := internal.ConnectRabbitMQ(constant.UsernameRabbitMQ, constant.PasswordRabbitMQ, constant.URLRabbitMQ, constant.VhostRabbitMQ, constant.CertPem, constant.ClientCertPem, constant.ClientKeyPem)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	consumeConn, err := internal.ConnectRabbitMQ(constant.UsernameRabbitMQ, constant.PasswordRabbitMQ, constant.URLRabbitMQ, constant.VhostRabbitMQ, constant.CertPem, constant.ClientCertPem, constant.ClientKeyPem)
	if err != nil {
		panic(err)
	}
	defer consumeConn.Close()

	consumeClient, err := internal.NewRabbitMQClient(consumeConn)
	if err != nil {
		panic(err)
	}
	defer consumeClient.Close()

	// Create Unnamed Queue which will generate a random name, set AutoDelete to True
	queue, err := consumeClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err := consumeClient.CreateBinding(queue.Name, queue.Name, "customer_callbacks"); err != nil {
		panic(err)
	}

	messageBus, err := consumeClient.Consume(queue.Name, "customer-api", true)
	if err != nil {
		panic(err)
	}
	go func() {
		for message := range messageBus {
			fmt.Printf("Message Callback %s\n", message.CorrelationId)
		}
	}()
	// Create context to manage timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create customer from sweden
	for i := 0; i < 10; i++ {
		if err := client.Send(ctx, "customers_events", "customers.created.se", amqp091.Publishing{
			ContentType:  "text/plain",       // The payload we send is plaintext, could be JSON or others..
			DeliveryMode: amqp091.Persistent, // This tells rabbitMQ that this message should be Saved if no resources accepts it before a restart (durable)
			Body:         []byte("An cool message between services"),
			// We add a REPLYTO which defines the
			ReplyTo: queue.Name,
			// CorrelationId can be used to know which Event this relates to
			CorrelationId: fmt.Sprintf("customer_created_%d", i),
		}); err != nil {
			panic(err)
		}
	}
	var blocking chan struct{}

	fmt.Println("Waiting on Callbacks, to close the program press CTRL+C")
	// This will block forever
	<-blocking
}
