package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jfelipeforero/iparking/internal"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
        conn, err := internal.ConnectRabbitMQ("user", "password", "localhost:5672", "booking")
        if err != nil {
                panic(err)
        }
        defer conn.Close()

        consumeConn, err := internal.ConnectRabbitMQ("user", "password", "localhost:5672", "booking")
        if err != nil {
                panic(err)
        }
        defer consumeConn.Close()

        client, err := internal.NewRabbitMQClient(conn)
        if err != nil {
                panic(err)
        }
        defer client.Close()

        consumeClient, err := internal.NewRabbitMQClient(consumeConn)
        if err != nil {
                panic(err)
        }
        defer consumeClient.Close()

        queue, err := consumeClient.CreateQueue("", true, true)
        if err != nil {
                panic(err)
        }

        if err := consumeClient.CreateBinding(queue.Name, queue.Name, "booking_callbacks"); err != nil {
                panic(err)
        }

        messageBus, err := consumeClient.Consume(queue.Name, "booking-api", true)
        if err != nil {
                panic(err)
        }

        go func() {
                for message := range messageBus {
                        log.Printf("Message callback: %s", message.CorrelationId)
                }
        }()

        ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
        defer cancel()

        for i := 0; i < 10; i++ {
                if err := client.Send(ctx, "booking_events", "booking.created.us", amqp091.Publishing{
                        ContentType: "text/plain",
                        DeliveryMode: amqp091.Persistent,
                        ReplyTo: queue.Name,
                        CorrelationId: fmt.Sprintf("customer_created_%d", i),
                        Body: []byte(`A cool message between services`),
                }); err != nil {
                        panic(err)
                }
        }

        // time.Sleep(10 * time.Second)

        log.Println(client)
        
        var blocking chan struct {}
        <-blocking

}
