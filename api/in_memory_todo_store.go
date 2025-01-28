package main

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{map[string]string{}}
}

type InMemoryTodoStore struct {
	store map[string]string
}

func (i *InMemoryTodoStore) CreateTodo(todoId string) {
	i.store[todoId] = "item"
}

func (i *InMemoryTodoStore) GetTodo(todoId string) string {
	return i.store[todoId]
}
