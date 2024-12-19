package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection // tcp connection used for all app
	ch   *amqp.Channel    // mutiplex connection, should be used for any new concurrent connection
}

// ConnectRabbitMQ will spawn a Connection
func ConnectRabbitMQ(username, password, host, vhost, caCert, clientCert, clientKey string) (*amqp.Connection, error) {
	ca, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}
	// Load the key pair
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	// Add the CA to the cert pool
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(ca)

	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	}
	// Setup the Connection to RabbitMQ host using AMQPs and Apply TLS config
	conn, err := amqp.DialTLS(fmt.Sprintf("amqps://%s:%s@%s/%s", username, password, host, vhost), tlsConf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	if err := ch.Confirm(false); err != nil {
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

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) (amqp.Queue, error) {
	q, err := rc.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return q, err
}

func (rc RabbitClient) CreateExchange(name, exchangeType string, durable, autoDelete bool) error {
	return rc.ch.ExchangeDeclare(name, exchangeType, durable, autoDelete, false, false, nil)
}

func (rc RabbitClient) DeleteExchange(name string) error {
	return rc.ch.ExchangeDelete(name, false, false)
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, true, false, options)
	if err != nil {
		return err
	}

	fmt.Println(confirmation.Wait())

	return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

/*
ApplyQos is used to apply qouality of service to the channel.

Prefetch count - How many messages the server will try to keep on the Channel.

Prefetch Size - How many Bytes the server will try to keep on the channel.

Global -- Any other Consumers on the connection in the future will apply the same rules if TRUE.
*/
func (rc RabbitClient) ApplyQos(count, size int, global bool) error {
	// Apply Quality of Serivce
	return rc.ch.Qos(
		count,
		size,
		global,
	)
}
