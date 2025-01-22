package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// declare user input variable
	var task string
	var taskList []string

	// prompt user for input
	fmt.Println("Enter a task:")

	// capture user task
	input := bufio.NewReader(os.Stdin)
	task, err := input.ReadString('\n')

	// trim the endline '\n' character
	task = task[:len(task)-1]
	if err != nil {
		fmt.Println(err)
	} else if len(task) == 0 {
		fmt.Println("You forgot to write a task!")
	}

	if len(task) > 0 {
		// add to list of tasks
		taskList = append(taskList, task)
		// print the task to console
		fmt.Println("Your task: ", task)
	}

	if len(taskList) > 0 {
		fmt.Println(taskList)
	} else {
		fmt.Println("You don't have any tasks!")
	}

}

// flags package for reading off the command line