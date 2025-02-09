package counter

import (
	"os"
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	testContent := "This is a test file. This file has some text. This text is for testing purposes."

	filename := "testfile.txt"
	err := os.WriteFile(filename, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Error creating the file: %v", err)
	}
	defer os.Remove(filename)

	wordCounts, err := CountWords(filename)
	if err != nil {
		t.Fatalf("CountWords failed: %v", err)
	}

	expectedWordCounts := map[string]int{
		"this":     3,
		"is":       2,
		"a":        1,
		"test":     1,
		"file":     2,
		"has":      1,
		"some":     1,
		"text":     2,
		"for":      1,
		"testing":  1,
		"purposes": 1,
	}

	if !reflect.DeepEqual(wordCounts, expectedWordCounts) {
		t.Errorf("CountWords failed: expected: %v, got: %v", expectedWordCounts, wordCounts)
	}
}
