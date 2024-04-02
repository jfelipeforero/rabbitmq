package main

import (
	"context"
	"log"
	"time"

	"github.com/jfelipeforero/iparking/internal"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func main() { 
        conn, err := internal.ConnectRabbitMQ("user", "password", "localhost:5672", "booking")
        if err != nil {
                panic(err)
        }
        defer conn.Close()

        publishConn, err := internal.ConnectRabbitMQ("user", "password", "localhost:5672", "booking")
        if err != nil {
                panic(err)
        }
        defer conn.Close()

        client, err := internal.NewRabbitMQClient(conn)
        if err != nil {
                panic(err)
        }
        defer client.Close()

        publishClient, err := internal.NewRabbitMQClient(publishConn)
        if err != nil {
                panic(err)
        }
        defer publishClient.Close()

        queue, err := client.CreateQueue("", true, true)
        if err != nil {
                panic(err)
        }

        if err := client.CreateBinding(queue.Name, "", "booking_events"); err != nil {
                panic(err)
        }

        messageBus, err := client.Consume(queue.Name, "email-service", false)
        if err != nil {
                panic(err)
        }

        var blocking chan struct{}

        ctx := context.Background()

        ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
        defer cancel()

        errg, ctx := errgroup.WithContext(ctx)

        if err := client.ApplyQoS(10, 0, true); err != nil {
                panic(err)
        }

        errg.SetLimit(10)

        go func() {
                for message := range messageBus {
                        msg := message
                        errg.Go(func() error {
                                log.Printf("New message: %v", msg)
                                time.Sleep(10 * time.Second)
                                if err := msg.Ack(false); err != nil {
                                        log.Println("Ack message failed")
                                        return err
                                } 
                                if err := publishClient.Send(ctx, "booking_callbacks", msg.ReplyTo, amqp091.Publishing{
                                        ContentType: "text/plain",
                                        DeliveryMode: amqp091.Persistent,
                                        Body: []byte("RPC Complete!"),
                                        CorrelationId: msg.CorrelationId,
                                }); err != nil {
                                        panic(err)
        
                                }
                                log.Println("Acknowledge message: %s\n", message.MessageId)
                                return nil
                        })
                }
        }()

        log.Println("Consuming events")
        <-blocking
}
