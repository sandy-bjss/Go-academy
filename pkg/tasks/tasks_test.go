package tasks_test

import (
	"encoding/json"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"sandy.goacademy/taskmaster/pkg/tasks"
)

func TestGetTasks(t *testing.T) {

	testTasks := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}, {Id: "02", Status: "started", Item: "Task 2"}}

	// setup temp file for testing
	tmpTaskFileName, teardown := setupTempTaskFile(t, testTasks)
	defer teardown()

	got := tasks.GetTasks(tmpTaskFileName)
	want := testTasks
	assertTaskListsEqual(t, got, want)

}

func TestCreateTasks(t *testing.T) {

	testTasks := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}}

	// setup temp file for testing
	tmpTaskFileName, teardown := setupTempTaskFile(t, testTasks)
	defer teardown()

	got := tasks.CreateTask(tasks.Task{Id: "02", Status: "started", Item: "Task 2"}, tmpTaskFileName)
	want := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}, {Id: "02", Status: "started", Item: "Task 2"}}
	assertTaskListsEqual(t, got, want)
}

func TestUpdateTasks(t *testing.T) {

	testTasks := []tasks.Task{{Id: "01", Status: "starting", Item: "Task 1"}}

	// setup temp file for testing
	tmpTaskFileName, teardown := setupTempTaskFile(t, testTasks)
	defer teardown()

	got := tasks.UpdateTask(tasks.Task{Id: "01", Status: "complete", Item: "Task 1"}, tmpTaskFileName)
	want := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}}
	assertTaskListsEqual(t, got, want)
}

func TestDeleteTasks(t *testing.T) {

	testTasks := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}, {Id: "02", Status: "started", Item: "Task 2"}}

	// setup temp file for testing
	tmpTaskFileName, teardown := setupTempTaskFile(t, testTasks)
	defer teardown()

	got := tasks.DeleteTask("02", tmpTaskFileName)
	want := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}}
	assertTaskListsEqual(t, got, want)
}

func assertTaskListsEqual(t testing.TB, got, want []tasks.Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func setupTempTaskFile(t testing.TB, content []tasks.Task) (string, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "task_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	// write contents to file
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		slog.Error("Could not marshal test tasks")
	}
	if _, err := tmpfile.Write(jsonBytes); err != nil {
		t.Fatalf("Filed to write content to temp file: %v", err)
	}

	// close file
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// return tempfile name and teardown func
	teardown := func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Fatalf("failed to remove temp file: %v", err)
		}
	}

	return tmpfile.Name(), teardown
}
