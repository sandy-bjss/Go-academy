package main

import (
	"fmt"
	"os"
	"os/signal"
)

// ===================================================================
// MAIN LAYER
// ===================================================================
func main() {
	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)
	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Starting server...\nCTRL-C to shutdown")

	// start api
	go Api()

	// capture signal
	<-sigs
	fmt.Println("\nKeyboard interupt...\nShutting down server")
}

//===================================================================

// ===================================================================
// API LAYER
// ===================================================================
func api() {
	// start server
	// define endpoints
}
// define handlers
//===================================================================


// ===================================================================
// SERVICE LAYER
// ===================================================================
func actor() {

}
//===================================================================

// ===================================================================
// CRUD LAYER
// ===================================================================
func loadData() {
	//load data
}

func updateData() {
	//update data
}

func deleteData() {
	//delete data
}
//===================================================================