package main

import (
	"fmt"
	"net/http"
	"strings"
)

type TodoStore interface {
	GetTodo(todo string) string
	CreateTodo(todoId string)
}

type TodoServer struct {
	store TodoStore
}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		t.ProcessTodo(w, r)
	case http.MethodGet:
		t.ShowTodo(w, r)
	}
}

func (t *TodoServer) ProcessTodo(w http.ResponseWriter, r *http.Request) {
	todoId := strings.TrimPrefix(r.URL.Path, "/todos/")
	t.store.CreateTodo(todoId)
	w.WriteHeader(http.StatusAccepted)

}

func (t *TodoServer) ShowTodo(w http.ResponseWriter, r *http.Request) {
	todoId := strings.TrimPrefix(r.URL.Path, "/todos/")

	todo := t.store.GetTodo(todoId)

	if todo == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, todo)
}
