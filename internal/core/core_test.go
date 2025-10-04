package core_test

import (
	"fizzbuzz/internal/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	// Create a new Core instance for tests
	core := core.NewCore()

	tests := []struct {
		name     string
		numbers  []int
		values   []string
		limit    int
		expected string
		err      error
	}{
		{
			name:     "FizzBuzz case",
			numbers:  []int{3, 5},
			values:   []string{"Fizz", "Buzz"},
			limit:    5,
			expected: "[1 2 Fizz 4 Buzz]",
			err:      nil,
		},
		{
			name:     "Single number case",
			numbers:  []int{2},
			values:   []string{"Even"},
			limit:    4,
			expected: "[1 Even 3 Even]",
			err:      nil,
		},
		{
			name:     "Empty numbers and values",
			numbers:  []int{},
			values:   []string{},
			limit:    3,
			expected: "[1 2 3]",
			err:      nil,
		},
		{
			name:     "Limit zero",
			numbers:  []int{3},
			values:   []string{"Fizz"},
			limit:    0,
			expected: "[]",
			err:      nil,
		},
		{
			name:     "Multiple matches",
			numbers:  []int{2, 4},
			values:   []string{"Two", "Four"},
			limit:    4,
			expected: "[1 Two 3 TwoFour]",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := core.ParseMessage(tt.numbers, tt.values, tt.limit)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
