package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTodosAndGetThem(t *testing.T) {
	store := NewInMemoryTodoStore()
	server := TodoServer{store}
	todo := "04"

	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), newPostTodoRequest(todo))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetTodoRequest(todo))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "item")
}
