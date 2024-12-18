package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lechuhuuha/eda-rabbitmq/constant"
	"github.com/lechuhuuha/eda-rabbitmq/internal"
	"golang.org/x/sync/errgroup"
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

	q, err := client.CreateQueue("", true, false)
	if err != nil {
		panic(err)
	}

	if err := client.CreateBinding(q.Name, "", constant.ExchangeEvent); err != nil {
		panic(err)
	}

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blockCh chan struct{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				fmt.Println("New message", msg)
				if err := msg.Ack(false); err != nil {
					fmt.Println("Ack mess failed", err)
					return err
				}
				fmt.Println("Acked message")
				return nil
			})
		}
	}()

	<-blockCh
}
