package main

import (
	"encoding/json"
	"flag"
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

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func printTasks(tasks Tasks) {
	// iterate through the tasks and print
	for i := 0; i < len(tasks.Tasks); i++ {
		fmt.Println("Task id: " + tasks.Tasks[i].Id)
		fmt.Println("Task status: " + tasks.Tasks[i].Status)
		fmt.Println("Task item: " + tasks.Tasks[i].Item)
	}
}

func main() {

	// TODO

	// read file test.txt
	f, err := os.Open("tasks.json")
	check(err)
	defer f.Close()

	byteArray, err := io.ReadAll(f)
	check(err)

	// initialise tasks variable
	var tasks Tasks

	// unmarhsal byte array
	json.Unmarshal(byteArray, &tasks)

	// capture cli flags
	// define flags, note: return pointers not values
	actionPtr := flag.String("do", "a", "specify what action to perform; 'a' to add, 'u' to update or 'd' to delete")
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
	}

	// print the tasks
	printTasks(tasks)

	// save the new tasks list to a json file
	jsonBytes, err := json.Marshal(tasks)
	check(err)

	err = os.WriteFile("tasks.json", jsonBytes, 0644)
}
