package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	// connect to rabbitmq
	rabbitconn, err := connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitconn.Close()

	app := Config{Rabbit: rabbitconn}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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