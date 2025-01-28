package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubTodoStore struct {
	todos     map[string]string
	todoCalls []string
}

func (s *StubTodoStore) GetTodo(todoId string) string {
	if todoId == "" {
		todoIds := fmt.Sprint(s.todos)
		return todoIds
	}
	todo := s.todos[todoId]
	return todo
}

func (s *StubTodoStore) CreateTodo(todoId string) {
	s.todoCalls = append(s.todoCalls, todoId)
}

func TestGETTodo(t *testing.T) {
	store := StubTodoStore{
		map[string]string{
			"01": "todo status 1 item 1",
			"02": "todo stauts 2 item 2",
		},
		nil,
	}

	server := &TodoServer{&store}

	t.Run("returns all todo list ids", func(t *testing.T) {
		request := newGetTodoRequest("")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "map[01:todo status 1 item 1 02:todo stauts 2 item 2]"

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})
	t.Run("return first item in todo list", func(t *testing.T) {
		request := newGetTodoRequest("01")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "todo status 1 item 1"

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})
	t.Run("returns 404 on missing todo ID", func(t *testing.T) {
		request := newGetTodoRequest("missingID")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatus(t, got, want)
	})
}

func TestCreateTodo(t *testing.T) {
	store := StubTodoStore{map[string]string{}, nil}

	server := &TodoServer{&store}

	t.Run("it return acceptance on POST", func(t *testing.T) {
		todoId := "03"
		request := newPostTodoRequest(todoId)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.todoCalls) != 1 {
			t.Errorf("got %d calls to CreatTodo, want %d", len(store.todoCalls), 1)
		}
		if store.todoCalls[0] != todoId {
			t.Errorf("did not creat correct todo, got %s want %s", store.todoCalls[0], todoId)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetTodoRequest(todoId string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/todos/%s", todoId), nil)
	return req
}

func newPostTodoRequest(todoId string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/todos/%s", todoId), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}
