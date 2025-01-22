package main

import (
	"encoding/json"
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

func main() {

	// TODO
	
	// read file test.txt
	f, err := os.Open("tasks.json")
	check(err)

	byteArray, err := io.ReadAll(f)
	check(err)

	// initialise tasks variable
	var tasks Tasks

	// unmarhsal byte array
	json.Unmarshal(byteArray, &tasks)

	// iterate through the tasks and print
	for i := 0; i < len(tasks.Tasks); i++ {
		fmt.Println("Task id: " + tasks.Tasks[i].Id)
		fmt.Println("Task status: " + tasks.Tasks[i].Status)
		fmt.Println("Task item: " + tasks.Tasks[i].Item)
	}

	f.Close()


	
	// use JSON filetype

}
