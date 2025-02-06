package tasks_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"sandy.goacademy/taskmaster/pkg/tasks"
)

func TestGetTasks(t *testing.T) {

	jsonContent := `[
		{"id: "01", "status": "complete", "item": "Task 1" },
		{"id: "02", "status": "started", "item": "Task 2" }
	]`

	// setup temp file for testing
	tmpTaskFileName, teardown := setupTempTaskFile(t, jsonContent)
	defer teardown()

	got := tasks.GetTasks(tmpTaskFileName)

	want := []tasks.Task{{Id: "01", Status: "complete", Item: "Task 1"}, {Id: "02", Status: "started", Item: "Task 2"}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}

}

func setupTempTaskFile(t testing.TB, content string) (string, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "task_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	fmt.Println(tmpfile.Name())

	// write contents to file
	if _, err := tmpfile.Write([]byte(content)); err != nil {
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
