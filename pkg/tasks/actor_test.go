package tasks_test

import (
	"fmt"
	"testing"
)

func TestParallel(t *testing.T) {
	t.Run("Subtest 1", func(t *testing.T) {
		t.Parallel()
		fmt.Println("Hello parallel test 1")
	})

	t.Run("Subtest 2", func(t *testing.T) {
		t.Parallel()
		fmt.Println("Hello parallel test 2")
	})

	t.Run("Subtest 3", func(t *testing.T) {
		t.Parallel()
		fmt.Println("Hello parallel test 3")
	})

	t.Run("Subtest 4", func(t *testing.T) {
		t.Parallel()
		fmt.Println("Hello parallel test 4")
	})
}
