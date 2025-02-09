package greeting

import (
	"testing"
)

func TestSayHello(t *testing.T) {
	testName := "Test User"
	expected := "Hello, Test User!\n"

	// This will capture any output that is printed in the console
	output := captureOutput(func() {
		SayHello(testName)
	})

	if output != expected {
		t.Errorf("SayHello() failed, expected %q, got %q", expected, output)
	}
}
