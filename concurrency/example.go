package main

import (
	"fmt"
	"time"
)

func Example() {
	fmt.Println("Hello, concurrency!\n===================")

	ch := make(chan string)

	// produce data
	go sendData(ch)

	// concurrently consume data from channel
	for i := 1; i <= 4; i++ {
		go receiveData(ch)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("===================\nBye, concurrency!")
}

func sendData(c chan string) {
	c <- "Hello"
	c <- "Channel"
	c <- "Based"
	c <- "World"
}

func receiveData(c chan string) {
	chData := <-c
	fmt.Println(chData)
}
