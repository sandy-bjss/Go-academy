package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"sandy.goacademy/taskmaster/pkg/tasks"
)

const TASK_LIST_JSON_FILE = "../../tasks.json"
const TraceIDString = TraceIDType("traceID")
const PORT = ":8080"

func main() {
	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)
	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Starting server...\nCTRL-C to shutdown")

	// start api
	go Api()

	// capture signal
	<-sigs
	fmt.Println("\nKeyboard interupt...\nShutting down server")
}

type TraceIDType string

func Api() {
	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// define handlers for routing
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

func middlewareTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("middlewareTRaceID: adding context traceID")

		ctx := context.WithValue(r.Context(), string(TraceIDString), uuid.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing GET handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	taskList := tasks.LoadTasks(TASK_LIST_JSON_FILE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing CREATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	var newTask tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		fmt.Printf("server: invlaid JSON: %s\n", err)
		slog.Error("Invlaid JSON, could not create new task")
		return
	}

	taskList := tasks.LoadTasks(TASK_LIST_JSON_FILE)
	taskList = tasks.AddTask(newTask, taskList)
	tasks.SaveTasks(taskList, TASK_LIST_JSON_FILE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing UPDATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	var updatedTask tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		fmt.Printf("server: invlaid JSON: %s\n", err)
		slog.Error("Invlaid JSON, could not create new task")
		return
	}
	if updatedTask.Id == "" {
		slog.Error("No task ID supplied, can't update task list")
		return
	}

	taskList := tasks.LoadTasks(TASK_LIST_JSON_FILE)
	taskList = tasks.UpdateTask(updatedTask, taskList)
	tasks.SaveTasks(taskList, TASK_LIST_JSON_FILE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing DELETE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	taskIdToDelete := r.URL.Query().Get("id")
	if taskIdToDelete == "" {
		slog.Error("No Task ID supplied, nothing deleted.")
	}

	taskList := tasks.LoadTasks(TASK_LIST_JSON_FILE)
	taskList = tasks.DeleteTask(taskIdToDelete, taskList)
	tasks.SaveTasks(taskList, TASK_LIST_JSON_FILE)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}
