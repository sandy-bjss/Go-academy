package main

import (
	"fmt"
	"time"
)

type operation struct {
	action    string
	parameter string
	response  chan string
}

var requests chan operation = make(chan operation)
var done chan struct{} = make(chan struct{})

func main() {
	fmt.Println("Hello, please place your order :)\n===================")

	Start()
	defer Stop()

	go OrderCoffee("Flat white")
	go OrderCoffee("Double esspresso")
	go fmt.Println("Last made coffee: " + GetLastCoffeeMade())
	go OrderCoffee("Aeropress")
	go fmt.Println("Last made coffee: " + GetLastCoffeeMade())
	go OrderCoffee("Peppermint tea")
	go fmt.Println("Last made coffee: " + GetLastCoffeeMade())

	// don't need sleep now Stop function implemented
	// main won't complete until Stop() function has completed
	// time.Sleep(10 * time.Second)
}

func OrderCoffee(coffeeOrder string) {
	order := operation{action: "order", parameter: coffeeOrder}

	requests <- order
}

func Start() {
	go monitorRequests()
}

func Stop() {
	shutdown := operation{action: "shutdown", parameter: "", response: nil}
	requests <- shutdown
	<-done
}

func monitorRequests() {
	// protected data
	var lastCoffeeMade string

	for op := range requests {
		switch op.action {
		case "order":
			requestCoffee := op.parameter
			makeCoffee(requestCoffee)
			lastCoffeeMade = requestCoffee
		case "lastmade":
			op.response <- lastCoffeeMade
		case "shutdown":
			// stop accepting new requests
			fmt.Println("===================\nSorry, coffe shop is closing. We can't accept any new orders.")
			close(requests)
		}
	}
	// signal all requests completed
	fmt.Println("All coffees have been made. We are now closed!")
	close(done)
}

func makeCoffee(coffee string) {
	fmt.Println("Brewing " + coffee)
	time.Sleep(2 * time.Second)
	fmt.Println(coffee + " now ready")
}

func GetLastCoffeeMade() string {
	answer := make(chan string)

	op := operation{
		action:    "lastmade",
		parameter: "",
		response:  answer,
	}

	requests <- op
	return <-answer
}
