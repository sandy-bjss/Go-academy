package main

import (
	"fmt"
	"os"
	"os/signal"
)

const PORT = ":8080"

func main() {
	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)
	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Starting server...\nCTRL-C to shutdown")

	// start api
	go FileServer()

	// capture signal
	<-sigs
	fmt.Println("\nKeyboard interupt...\nShutting down server")
}
