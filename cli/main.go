package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {

	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)

	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)

	// run cli
	TodoCli()
	fmt.Println("==================\ncli running: CTRL-C to exit")

	// capture signal
	sig := <-sigs
	fmt.Println("\nKeyboard interupt: " + sig.String() + "\nClosing todo app")
}
