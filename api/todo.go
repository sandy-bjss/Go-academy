package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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

// Add task method
func (ts *Tasks) AddTask(addTask Task) {
	ts.Tasks = append(ts.Tasks, addTask)
}

// Update task method
func (ts *Tasks) UpdateTask(updateTask Task) error {
	for idx, t := range ts.Tasks {
		if t.Id == updateTask.Id {
			// need to access the task via task[idx] to modify value in memory
			ts.Tasks[idx].Status = updateTask.Status
			ts.Tasks[idx].Item = updateTask.Item
			return nil
		}
	}
	return errors.New("uh oh: task not found")
}

// Delete task method
func (ts *Tasks) DeleteTask(deleteTask Task) error {
	for idx, t := range ts.Tasks {
		if t.Id == deleteTask.Id {
			ts.Tasks = append(ts.Tasks[:idx], ts.Tasks[idx+1:]...)
			return nil
		}
	}
	return errors.New("Uh oh: did not find task, so could perform delete")
}

func printTasks(tasks Tasks) {
	// iterate through the tasks and print
	for i := 0; i < len(tasks.Tasks); i++ {
		fmt.Println("Task id: " + tasks.Tasks[i].Id)
		fmt.Println("Task status: " + tasks.Tasks[i].Status)
		fmt.Println("Task item: " + tasks.Tasks[i].Item)
	}
}

func loadTodos() {
	// read todo JSON data
	f, err := os.Open("tasks.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	byteArray, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	// initialise tasks variable
	var tasks Tasks

	// unmarhsal byte array
	json.Unmarshal(byteArray, &tasks)

	// print the tasks
	printTasks(tasks)
}
