package main

import (
	"fmt"
	"net/http"
	"strings"
)

type TodoStore interface {
	GetTodos(todo string) string
}

type TodoServer struct {
	store TodoStore
}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoId := strings.TrimPrefix(r.URL.Path, "/todos/")

	todo := t.store.GetTodos(todoId)

	if todo == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, todo)
}
