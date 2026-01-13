package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	ConnString := "amqp://guest:guest@localhost:5672/"
	cxn, err := amqp.Dial(ConnString)
	defer cxn.Close()
	if err == nil {
		fmt.Println("Connection to server successful!")
	}

	// server shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nPeril server shutting down...")
}
