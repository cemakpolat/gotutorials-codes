package calculator

import "testing"

func TestAdd(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{10, 5, 15},
		{0, 0, 0},
		{-5, 5, 0},
		{-10, -5, -15},
	}
	for _, tc := range testCases {
		if actual := Add(tc.a, tc.b); actual != tc.expected {
			t.Errorf("Add(%d,%d) failed, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestMultiply(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{10, 5, 50},
		{0, 0, 0},
		{-5, 5, -25},
		{-10, -5, 50},
	}
	for _, tc := range testCases {
		if actual := Multiply(tc.a, tc.b); actual != tc.expected {
			t.Errorf("Multiply(%d,%d) failed, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
}
