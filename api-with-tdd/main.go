package main

import (
	"log"
	"net/http"
)

func main() {
	server := &TodoServer{NewInMemoryTodoStore()}
	log.Fatal((http.ListenAndServe(":5000", server)))
}
