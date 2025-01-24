package todo

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

// Tasks struct which contains an array of tasks
type Tasks struct {
	Tasks []Task `json:"Tasks"`
}

// Task struct which contains a status and an item.
type Task struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Item   string `json:"item"`
}

// const for traceID key of TraceIDType
const TraceIDString = TraceIDType("traceID")

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

func printTasks(tasks Tasks) {
	// iterate through the tasks and print
	for i := 0; i < len(tasks.Tasks); i++ {
		fmt.Println("Task id: " + tasks.Tasks[i].Id)
		fmt.Println("Task status: " + tasks.Tasks[i].Status)
		fmt.Println("Task item: " + tasks.Tasks[i].Item)
	}
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

	// read file test.txt
	f, err := os.Open("tasks.json")
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
	var tasks Tasks

	// unmarhsal byte array
	json.Unmarshal(byteArray, &tasks)

	// capture cli flags
	// define flags, note: return pointers not values
	actionPtr := flag.String("do", "", "specify what action to perform; 'a' to add, 'u' to update or 'd' to delete")
	idPtr := flag.String("id", "foo", "task id")
	statusPtr := flag.String("status", "bar", "task status")
	itemPtr := flag.String("item", "foobar", "task item")
	flag.Parse()

	// check what to do with the supplied item
	switch *actionPtr {
	case "a":
		// create task item from flags
		var newTask Task
		// set task values
		newTask.Id = *idPtr
		newTask.Status = *statusPtr
		newTask.Item = *itemPtr
		// add this task to tasks array
		tasks.Tasks = append(tasks.Tasks, newTask)
	case "u":
		fmt.Println("inside update tasks")
		// loop through tasks
		for idx, t := range tasks.Tasks {
			if t.Id == *idPtr {
				// need to access the task via task[idx] to modify value in memory
				tasks.Tasks[idx].Status = *statusPtr
				tasks.Tasks[idx].Item = *itemPtr
				break
			}
		}
	case "d":
		fmt.Println("inside delete tasks")
		// loop through tasks
		for idx, t := range tasks.Tasks {
			if t.Id == *idPtr {
				tasks.Tasks = append(tasks.Tasks[:idx], tasks.Tasks[idx+1:]...)
			}
		}
	default:
		break
	}

	// print the tasks
	printTasks(tasks)

	// save the new tasks list to a json file
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println(err)
		logger.ErrorContext(ctx, "Encountered an Error", "error", err)
	}

	err = os.WriteFile("tasks.json", jsonBytes, 0644)
	if err == nil {
		slog.InfoContext(ctx, "Saved tasks to file 'tasks.json'")
	}
}
