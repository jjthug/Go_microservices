package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/events"
	"log"
	"math"
	"os"
	"time"
)

func main() {

	// connect to rabbitmq
	rabbitconn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitconn.Close()
	// listen for messages

	log.Println("listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := events.NewConsumer(rabbitconn)
	if err != nil {
		panic(err)
	}

	// consume the messages
	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready yet....")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			return connection, nil
		}
		if counts > 5 {
			log.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off....")
		time.Sleep(backOff)
		continue
	}
}