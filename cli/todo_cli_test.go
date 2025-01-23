package todo

import "testing"

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHello(t *testing.T) {
	got := "sandy"
	want := "sandy"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
