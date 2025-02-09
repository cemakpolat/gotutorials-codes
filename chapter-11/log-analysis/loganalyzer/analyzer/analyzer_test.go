package analyzer

import (
	"os"
	"reflect"
	"testing"
)

func TestAnalyzeLogs(t *testing.T) {
	testLog := `[info] this is an info log message
[error] this is an error message
[warning] this is a warning message
[info] another info message
[error] another error message with more details
[debug] some debug information
`

	filename := "test.log"

	err := os.WriteFile(filename, []byte(testLog), 0644)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer os.Remove(filename)

	t.Run("Test without filter", func(t *testing.T) {
		severityCounts, err := AnalyzeLogs(filename, "")
		if err != nil {
			t.Errorf("AnalyzeLogs failed with: %v", err)
		}

		expected := map[string]int{
			"info":    2,
			"error":   2,
			"warning": 1,
			"debug":   1,
		}
		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed: expected: %v, got: %v", expected, severityCounts)
		}
	})

	t.Run("Test with filter", func(t *testing.T) {
		severityCounts, err := AnalyzeLogs(filename, "error")
		if err != nil {
			t.Errorf("AnalyzeLogs failed with filter, error: %v", err)
		}
		expected := map[string]int{
			"error": 2,
		}

		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed with filter, expected: %v, got: %v", expected, severityCounts)
		}
	})

	t.Run("Test with invalid logs", func(t *testing.T) {
		testLog := `invalid log`
		err = os.WriteFile(filename, []byte(testLog), 0644)
		if err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
		severityCounts, err := AnalyzeLogs(filename, "")
		if err != nil {
			t.Fatalf("AnalyzeLogs should not have failed")
		}

		expected := map[string]int{}
		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed with invalid logs, expected %v got %v", expected, severityCounts)
		}
	})
}
