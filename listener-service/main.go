package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to conect to rabbiqmq
	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitCon.Close()

	// start listening to messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitCon)
	if err != nil {
		panic(err)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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
