package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// const for traceID key of TraceIDType
const TraceIDString = TraceIDType("traceID")

// Task struct which contains a status and an item.
type Task struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Item   string `json:"item"`
}

var tasks []Task
var tasksData = "tasks.json"

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

func loadTasks() {
	file, err := os.Open(tasksData)
	if err != nil {
		slog.Info("no json file with todos exists, a blank Todo slice has been initialised")
		tasks = []Task{}
	}
	defer file.Close()

	byteArray, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		slog.Error("Could not read file data, a blank Todo slice has been initialised")
		tasks = []Task{}
	}

	json.Unmarshal(byteArray, &tasks)
}

func saveTasks() {
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		slog.Error("Could not save tasks")
	}

	os.WriteFile("tasks.json", jsonBytes, 0644)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing GET handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	loadTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing CREATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		fmt.Printf("server: invlaid JSON: %s\n", err)
		slog.Error("Invlaid JSON, could not create new task")
		return
	}

	loadTasks()
	tasks = append(tasks, newTask)
	saveTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing UPDATE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		fmt.Printf("server: invlaid JSON: %s\n", err)
		slog.Error("Invlaid JSON, could not create new task")
		return
	}
	if updatedTask.Id == "" {
		slog.Error("No task ID supplied, can't update task list")
		return
	}

	loadTasks()
	for i, t := range tasks {
		if t.Id == updatedTask.Id {
			tasks[i].Status = updatedTask.Status
			tasks[i].Item = updatedTask.Item
			break
		}
	}
	saveTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing DELETE handler", string(TraceIDString), r.Context().Value(string(TraceIDString)))

	taskIdToDelete := r.URL.Query().Get("id")
	if taskIdToDelete == "" {
		slog.Error("No Task ID supplied, nothing deleted.")
	}

	loadTasks()

	for i, t := range tasks {
		if t.Id == taskIdToDelete {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
		if i == len(tasks) {
			slog.Error("Could not find task to delete")
		}
	}

	saveTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
