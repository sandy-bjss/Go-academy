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

func LoadTasks(taskJSONfilestore string) []Task {
	file, err := os.Open(taskJSONfilestore)
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

func SaveTasks(tasks []Task, taskJSONfilestore string) {
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		slog.Error("Could not save tasks")
	}

	os.WriteFile(taskJSONfilestore, jsonBytes, 0644)
}

func GetTasks(taskJSONfilestore string) []Task {
	return LoadTasks(taskJSONfilestore)
}

func CreateTask(task Task, taskJSONfilestore string) []Task {
	taskList := LoadTasks(taskJSONfilestore)
	taskList = append(taskList, task)
	SaveTasks(taskList, taskJSONfilestore)
	return taskList
}

func UpdateTask(taskToUpdate Task, taskJSONfilestore string) []Task {
	taskList := LoadTasks(taskJSONfilestore)
	for i, t := range taskList {
		if t.Id == taskToUpdate.Id {
			taskList[i].Status = taskToUpdate.Status
			taskList[i].Item = taskToUpdate.Item
			break
		}
	}
	SaveTasks(taskList, taskJSONfilestore)
	return taskList
}

func DeleteTask(taskIdToDelete string, taskJSONfilestore string) []Task {
	taskList := LoadTasks(taskJSONfilestore)
	for i, t := range taskList {
		if t.Id == taskIdToDelete {
			taskList = append(taskList[:i], taskList[i+1:]...)
			break
		}
		if i == len(taskList) {
			slog.Error("Could not find task to delete")
		}
	}
	SaveTasks(taskList, taskJSONfilestore)
	return taskList
}
