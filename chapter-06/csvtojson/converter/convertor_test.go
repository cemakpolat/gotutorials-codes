package converter

import (
	"os"
	"reflect"
	"testing"
)

func TestProcessCSV(t *testing.T) {
	testCases := []struct {
		csvPath    string
		expected   []Person
		shouldFail bool
	}{
		{
			csvPath: "testdata/valid.csv",
			expected: []Person{
				{Name: "Alice", Age: 30, City: "New York"},
				{Name: "Bob", Age: 25, City: "Los Angeles"},
				{Name: "Charlie", Age: 35, City: "Chicago"},
			},
			shouldFail: false,
		},
		{
			csvPath:    "testdata/invalid.csv",
			expected:   nil,
			shouldFail: true,
		},
	}

	// create testdata folder
	os.Mkdir("testdata", 0755)

	// create valid test data
	os.WriteFile("testdata/valid.csv", []byte("name,age,city\nAlice,30,New York\nBob,25,Los Angeles\nCharlie,35,Chicago"), 0644)

	// create invalid test data
	os.WriteFile("testdata/invalid.csv", []byte("name,age,city,extra\nAlice,30,New York,extra\nBob,25,Los Angeles,extra\nCharlie,35,Chicago,extra"), 0644)

	defer os.RemoveAll("testdata")
	for _, tc := range testCases {
		actual, err := ProcessCSV(tc.csvPath)
		if tc.shouldFail {
			if err == nil {
				t.Errorf("TestProcessCSV(%s) failed, should fail but did not", tc.csvPath)
			}
		} else {
			if err != nil {
				t.Errorf("TestProcessCSV(%s) failed with err: %v", tc.csvPath, err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("TestProcessCSV(%s) failed: expected %v, got %v", tc.csvPath, tc.expected, actual)
			}
		}
	}
}
