package main

import "fmt"

// create type to pass to the requests queue
type operation struct {
	action    string
	parameter string
	response  chan string
}

// make a requests channel to queue operations
var requests chan operation = make(chan operation)

// create a done channel for graceful shutdown
var done chan struct{} = make(chan struct{})

// actor start
func Start() {
	go monitorRequests()
}

// graceful shutdown
func Stop() {
	shutdown := operation{action: "shutdown", parameter: "", response: nil}
	requests <- shutdown
	<-done
}

// logic for processing requests
func monitorRequests() {
	// loop over reqeusts for sequential processing
	for req := range requests {
		fmt.Println(req.action)
	}
}
