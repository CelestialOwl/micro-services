package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func connectMongo() (*mongo.Client, error) {
	//create options
	clientoptions := options.Client().ApplyURI(mongoURL)
	clientoptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(context.TODO(), clientoptions)
	log.Println("connected to the mongo")
	if err != nil {
		log.Println("Error connecting to the mongodb", err)
		return nil, err
	}
	return c, nil
}

func main() {
	// connect to mongo
	mongoClient, err := connectMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := Config{
		Models: data.New(client),
	}

	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	//start web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Println("starting service on port ", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting rpc server on port " + rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
