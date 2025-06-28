package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {
	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitCon.Close()
	app := Config{
		RabbitMQ: rabbitCon,
	}

	log.Printf("Starting broker service on port %s", webPort)

	serve := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Rabbitmq not yet ready")
			counts++
		} else {
			log.Println("connected to rabbitmq!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("unable to connect", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off for now")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
