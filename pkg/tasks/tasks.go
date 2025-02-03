package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Task struct which contains a status and an item.
type Task struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Item   string `json:"item"`
}

var tasks []Task

func loadTasks(taskJSONFile string) {
	file, err := os.Open(taskJSONFile)
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

func saveTasks(tasks []Task) {
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		slog.Error("Could not save tasks")
	}

	os.WriteFile("tasks.json", jsonBytes, 0644)
}

func addTask(task Task, tasks []Task) {
	tasks = append(tasks, task)
}

func updateTask(taskToUpdate Task, tasks []Task) {
	for i, t := range tasks {
		if t.Id == taskToUpdate.Id {
			tasks[i].Status = taskToUpdate.Status
			tasks[i].Item = taskToUpdate.Item
			break
		}
	}
}

func deleteTask(taskIdToDelete string, tasks []Task) {
	for i, t := range tasks {
		if t.Id == taskIdToDelete {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
		if i == len(tasks) {
			slog.Error("Could not find task to delete")
		}
	}
}
