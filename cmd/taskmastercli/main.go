package main

import (
	"fmt"

	"sandy.goacademy/taskmaster/pkg/tasks"
)

func main() {
	var task = tasks.Task{Id: "00", Status: "Pending", Item: "First CLI task"}

	fmt.Println(task)
}
