package internal

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection // tcp connection used for all app
	ch   *amqp.Channel    // mutiplex connection, should be used for any new concurrent connection
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	return err
}

func (rc RabbitClient) CreateExchange(name, exchangeType string, durable, autoDelete bool) error {
	return rc.ch.ExchangeDeclare(name, exchangeType, durable, autoDelete, false, false, nil)
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}
