package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"sandy.goacademy/taskmaster/pkg/tasks"
)

// const for traceID key of TraceIDType
const TraceIDString = TraceIDType("traceID")

const TASK_LIST_JSON_FILE = "../../tasks.json"

type TraceIDType string

type customHandler struct {
	slog.Handler
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIDString).(string); ok {
		r.AddAttrs(slog.String(string(TraceIDString), traceID))
	}
	if traceID, ok := ctx.Value(TraceIDString).(uuid.UUID); ok {
		r.AddAttrs(slog.String(string(TraceIDString), traceID.String()))
	}
	return h.Handler.Handle(ctx, r)

}

func printTasks(taskList []tasks.Task) {
	// iterate through the tasks and print
	for i := 0; i < len(taskList); i++ {
		fmt.Println("Task id: " + taskList[i].Id)
		fmt.Println("Task status: " + taskList[i].Status)
		fmt.Println("Task item: " + taskList[i].Item)
	}
}

func main() {
	var task = tasks.Task{Id: "00", Status: "Pending", Item: "First CLI task"}

	fmt.Println(task)

	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)

	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)

	// run cli
	TodoCli()
	fmt.Println("==================\ncli running: CTRL-C to exit")

	// capture signal
	sig := <-sigs
	fmt.Println("\nKeyboard interupt: " + sig.String() + "\nClosing todo app")
}

func TodoCli() {

	// initialise context
	ctx, ctxDone := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, TraceIDString, uuid.New())
	defer ctxDone()

	// initialise logger
	baseHandler := slog.NewJSONHandler(os.Stdout, nil) // &slog.HandlerOptions{AddSource: true}
	handler := &customHandler{Handler: baseHandler}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.InfoContext(ctx, "starting")

	// read todo JSON data
	f, err := os.Open(TASK_LIST_JSON_FILE)
	if err != nil {
		fmt.Println(err)
		logger.ErrorContext(ctx, "Encountered an Error", "error", err)
	}
	defer f.Close()

	byteArray, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		logger.ErrorContext(ctx, "Encountered an Error", "error", err)
	}

	// initialise tasks variable
	var taskList []tasks.Task

	// unmarhsal byte array
	json.Unmarshal(byteArray, &taskList)

	// capture cli flags
	// define flags, note: return pointers not values
	actionPtr := flag.String("do", "", "specify what action to perform; 'a' to add, 'u' to update or 'd' to delete")
	idPtr := flag.String("id", "foo", "task id")
	statusPtr := flag.String("status", "bar", "task status")
	itemPtr := flag.String("item", "foobar", "task item")
	flag.Parse()

	// create task item from flags
	var newTask tasks.Task
	// set task values
	newTask.Id = *idPtr
	newTask.Status = *statusPtr
	newTask.Item = *itemPtr

	// check what to do with the supplied item
	switch *actionPtr {
	case "a":
		taskList = tasks.AddTask(newTask, taskList)
	case "u":
		taskList = tasks.UpdateTask(newTask, taskList)
	case "d":
		taskList = tasks.DeleteTask(newTask.Id, taskList)
	default:
		break
	}

	// print the tasks
	printTasks(taskList)

	// save the new tasks list to a json file
	jsonBytes, err := json.Marshal(taskList)
	if err != nil {
		fmt.Println(err)
		logger.ErrorContext(ctx, "Encountered an Error", "error", err)
	}

	err = os.WriteFile(TASK_LIST_JSON_FILE, jsonBytes, 0644)
	if err == nil {
		slog.InfoContext(ctx, "Saved tasks to file 'tasks.json'")
	}
}
