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

var taskList []Task

func LoadTasks(taskJSONFile string) []Task {
	file, err := os.Open(taskJSONFile)
	if err != nil {
		slog.Info("no json file with todos exists, a blank Todo slice has been initialised")
		taskList = []Task{}
	}
	defer file.Close()

	byteArray, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		slog.Error("Could not read file data, a blank Todo slice has been initialised")
		taskList = []Task{}
	}

	json.Unmarshal(byteArray, &taskList)

	return taskList
}

func SaveTasks(tasks []Task, file string) {
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		slog.Error("Could not save tasks")
	}

	os.WriteFile(file, jsonBytes, 0644)
}

func AddTask(task Task, tasks []Task) {
	tasks = append(tasks, task)
}

func UpdateTask(taskToUpdate Task, tasks []Task) {
	for i, t := range tasks {
		if t.Id == taskToUpdate.Id {
			tasks[i].Status = taskToUpdate.Status
			tasks[i].Item = taskToUpdate.Item
			break
		}
	}
}

func DeleteTask(taskIdToDelete string, tasks []Task) {
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
