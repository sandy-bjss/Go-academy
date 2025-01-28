package main

import "fmt"

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{map[string]string{}}
}

type InMemoryTodoStore struct {
	todos map[string]string
}

func (i *InMemoryTodoStore) CreateTodo(todoId string) {
	i.todos[todoId] = "item"
}

func (i *InMemoryTodoStore) GetTodo(todoId string) string {
	if todoId == "" {
		todoIds := fmt.Sprint(i.todos)
		return todoIds
	}
	return i.todos[todoId]
}
