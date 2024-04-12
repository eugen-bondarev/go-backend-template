package util

import (
	"slices"
	"testing"
)

func Test_Diff(t *testing.T) {
	type testCaseInt struct {
		a        []int
		b        []int
		expected []int
	}

	testCasesInt := []testCaseInt{
		{
			a:        []int{1, 2, 3},
			b:        []int{3, 4, 5, 1},
			expected: []int{2},
		},
		{
			a:        []int{42, 69, 128},
			b:        []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			expected: []int{42, 128, 69},
		},
	}

	for _, tc := range testCasesInt {
		res := Diff(tc.a, tc.b)
		slices.Sort(res)
		slices.Sort(tc.expected)
		if !slices.Equal(res, tc.expected) {
			t.Errorf("expected %v, got %v", tc.expected, res)
		}
	}

	type testCaseStr struct {
		a        []string
		b        []string
		expected []string
	}

	testCasesStr := []testCaseStr{
		{
			a:        []string{"hi", "hi", "lorem", "ipsum"},
			b:        []string{"hi", "hi", "lorem", "ipsum"},
			expected: []string{},
		},
		{
			a:        []string{"hi"},
			b:        []string{"lorem", "ipsum"},
			expected: []string{"hi"},
		},
		{
			a:        []string{"1", "2", "3"},
			b:        []string{"3", "4", "5"},
			expected: []string{"1", "2"},
		},
	}

	for _, tc := range testCasesStr {
		res := Diff(tc.a, tc.b)
		slices.Sort(res)
		slices.Sort(tc.expected)
		if !slices.Equal(res, tc.expected) {
			t.Errorf("expected %v, got %v", tc.expected, res)
		}
	}
}
