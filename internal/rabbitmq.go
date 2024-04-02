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
        // The connection used by the client
        conn *amqp.Connection //TCP Connection
        // Channel is used to process / Send messages
        ch *amqp.Channel // Multiplexed subconnection(?)
}

func ConnectRabbitMQ(username, password, host, vhost, caCert, clientCert, clientKey string) (*amqp.Connection, error) {
        ca, err := os.ReadFile(caCert)
        if err != nil {
                return nil, err
        }
        cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
        if err != nil {
                return nil, err
        }
        rootCAs := x509.NewCertPool()
        rootCAs.AppendCertsFromPEM(ca)

        tlsCfg := &tls.Config{
                RootCAs: rootCAs,
                Certificates: []tls.Certificate{cert},
        }
        return amqp.DialTLS(fmt.Sprintf("amqps://%s:%s@%s/%s", username, password, host, vhost), tlsCfg)
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
                conn:conn,
                ch: ch,
        }, nil
}

func (rc RabbitClient) Close() error {
        err := rc.ch.Close()
        if err != nil {
                return err
        }
        return nil
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
        q, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
        if err != nil {
                return amqp.Queue{}, nil
        }
        return q, err
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
        return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
        confirmation, err :=  rc.ch.PublishWithDeferredConfirmWithContext(ctx,
                exchange,
                routingKey,
                // Mandatory: True: Return upon failure.
                true,
                // Immediate
                false,
                options,
        )

        if err != nil {
                return err
        }
        confirmation.Wait()
        return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
        return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc RabbitClient) ApplyQoS(count, size int, global bool) error {
        return rc.ch.Qos(count, size, global)
}
