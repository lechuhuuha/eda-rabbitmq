package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"github.com/lechuhuuha/eda-rabbitmq/constant"
	"github.com/lechuhuuha/eda-rabbitmq/internal"
)

func main() {
	conn, err := internal.ConnectRabbitMQ(constant.UsernameRabbitMQ, constant.PasswordRabbitMQ, constant.URLRabbitMQ, constant.VhostRabbitMQ, constant.CertPem, constant.ClientCertPem, constant.ClientKeyPem)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mqClient, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	publishConn, err := internal.ConnectRabbitMQ(constant.UsernameRabbitMQ, constant.PasswordRabbitMQ, constant.URLRabbitMQ, constant.VhostRabbitMQ, constant.CertPem, constant.ClientCertPem, constant.ClientKeyPem)
	if err != nil {
		panic(err)
	}
	defer publishConn.Close()

	publishClient, err := internal.NewRabbitMQClient(publishConn)
	if err != nil {
		panic(err)
	}

	// Create Unnamed Queue which will generate a random name, set AutoDelete to True
	queue, err := mqClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	// Create binding between the customers_events exchange and the new Random Queue
	// Can skip Binding key since fanout will skip that rule
	if err := mqClient.CreateBinding(queue.Name, "", "customers_events"); err != nil {
		panic(err)
	}

	messageBus, err := mqClient.Consume(queue.Name, "email-service", false)
	if err != nil {
		panic(err)
	}

	// Set a timeout for 15 secs
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create an Errgroup to manage concurrecy
	g, ctx := errgroup.WithContext(ctx)
	// Set amount of concurrent tasks
	g.SetLimit(10)

	// Apply Qos to limit amount of messages to consume
	if err := mqClient.ApplyQos(10, 0, true); err != nil {
		panic(err)
	}

	go func() {
		for message := range messageBus {

			// Spawn a worker
			msg := message
			g.Go(func() error {
				// Multiple means that we acknowledge a batch of messages, leave false for now
				if err := msg.Ack(false); err != nil {
					fmt.Printf("Acknowledged message failed: Retry ? Handle manually %s\n", msg.MessageId)
					return err
				}
				fmt.Printf("Acknowledged message, replying to %s\n", msg.ReplyTo)

				// Use the msg.ReplyTo to send the message to the proper Queue
				if err := publishClient.Send(ctx, "customer_callbacks", msg.ReplyTo, amqp091.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp091.Transient,
					Body:          []byte("RFC Complete"),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					panic(err)
				}

				return nil
			})
		}
	}()

	var blockCh chan struct{}

	<-blockCh
}
