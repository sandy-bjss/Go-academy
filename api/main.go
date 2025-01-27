package main

import (
	"log"
	"net/http"
)

type InMemoryTodoStore struct{}

func (i *InMemoryTodoStore) GetTodos(todoId string) string {
	return "01"
}

func main() {
	server := &TodoServer{&InMemoryTodoStore{}}
	log.Fatal((http.ListenAndServe(":5000", server)))
}
