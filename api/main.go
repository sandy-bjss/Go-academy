package main

import (
	"log"
	"net/http"
)

type InMemoryTodoStore struct{}

func (i *InMemoryTodoStore) GetTodo(todoId string) string {
	return "01"
}

func (i *InMemoryTodoStore) CreateTodo(todoId string) {}

func main() {
	server := &TodoServer{&InMemoryTodoStore{}}
	log.Fatal((http.ListenAndServe(":5000", server)))
}
