package main

import "testing"

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAddTodo(t *testing.T) {
	tasks := Tasks{}
	task := Task{Id: "t01", Status: "started", Item: "write unit tests"}

	tasks.AddTask(task)

	if len(tasks.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks.Tasks))
	}
	if tasks.Tasks[0].Id != task.Id {
		t.Errorf("got %s want %s", tasks.Tasks[0].Id, task.Id)
	}
}

func TestUpdateTodo(t *testing.T) {
	tasks := Tasks{Tasks: []Task{{Id: "t01", Status: "started", Item: "write unit tests"}}}
	updatedTask := Task{Id: "t01", Status: "completed", Item: "write unit tests to test function, not implementation"}
	err := tasks.UpdateTask(updatedTask)

	if err != nil {
		t.Fatalf("task update failed: %v", err)
	}
	if tasks.Tasks[0].Status != updatedTask.Status {
		t.Errorf("got %s want %s", tasks.Tasks[0].Status, updatedTask.Status)
	}
	if tasks.Tasks[0].Item != updatedTask.Item {
		t.Errorf("got %s want %s", tasks.Tasks[0].Item, updatedTask.Item)
	}
}

func TestDeleteTodo(t *testing.T) {
	tasks := Tasks{Tasks: []Task{{Id: "t01", Status: "started", Item: "write unit tests"}, {Id: "t02", Status: "todo", Item: "publsih cli"}}}
	remainingTask := Task{Id: "t01", Status: "started", Item: "write unit tests"}
	deleteTask := Task{Id: "t02", Status: "todo", Item: "publsih cli"}
	err := tasks.DeleteTask(deleteTask)

	if err != nil {
		t.Fatalf("delete task failed: %v", err)
	}
	if len(tasks.Tasks) != 1 {
		t.Errorf("got %d want %d", len(tasks.Tasks), 1)
	}
	if tasks.Tasks[0] != remainingTask {
		t.Errorf("got %v want %v", tasks.Tasks[0], remainingTask)
	}
}
