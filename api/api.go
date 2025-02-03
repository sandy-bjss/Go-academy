package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

// const for traceID key of TraceIDType
const TraceIDString = TraceIDType("traceID")

// in memory todos
var todos = map[string]string{
	"started":   "item 1",
	"completed": "item 2",
	"pending":   "item 3",
}

type TraceIDType string

func Api() {
	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// define handlers for routing
	gh := http.HandlerFunc(genericHandler)
	mux.Handle("/", middlewareOne(middlewareTwo(middlewareTraceID(gh))))
	mux.Handle("POST /create", middlewareTraceID(http.HandlerFunc(createHandler)))
	mux.Handle("GET /get", middlewareTraceID(http.HandlerFunc(getHandler)))
	mux.Handle("POST /update", middlewareTraceID(http.HandlerFunc(updateHandler)))
	mux.Handle("DELETE /delete", middlewareTraceID(http.HandlerFunc(deleteHandler)))

	// start server
	if err := http.ListenAndServe(PORT, mux); err != nil {
		logger.Error("API Server couldn't start")
		return
	}
	slog.Info("Server Started, listening...", "PORT", PORT)
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("executing middleware 1")
		next.ServeHTTP(w, r)
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("executing middleware 2")
		next.ServeHTTP(w, r)
	})
}

func middlewareTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("middlewareTRaceID: adding context traceID")

		ctx := context.WithValue(r.Context(), string(TraceIDString), uuid.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing generic handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))
	w.Header().Set("Content-Type", "application/json")
	endpoint := strings.TrimPrefix(r.URL.Path, "")
	fmt.Fprintf(w, "Api endpoint: %s\nEndpoint type: %v", endpoint, w.Header().Get("Content-Type"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing GET handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	jsonData, err := json.Marshal(todos)
	if err != nil {
		slog.Error("error loading todas", string(TraceIDString), r.Context().Value(string(TraceIDString)))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing CREATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("server: could not read request body: %s\n", err)
		return
	}
	fmt.Printf("server: request body: %s\n", reqBody)

	todos["test"] = string(reqBody)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing UPDATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	todoID := r.Header.Get("id") //"pending"
	slog.Info(todoID)

	for k := range todos {
		if k == todoID {
			todos[k] = "updated to new item"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing DELETE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	taskToDelete := r.URL.Query().Get("id")
	if taskToDelete == "" {
		slog.Error("No Task ID supplied, nothing deleted.")
	}

	_, ok := todos[taskToDelete]
	if ok {
		delete(todos, taskToDelete)
	} else {
		slog.Error("No such task found in list of todos.")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
