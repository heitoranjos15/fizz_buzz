package core_test

import (
	"fizzbuzz/internal/core"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) SaveMessage(words []int, multiples []string, limit int) error {
	args := m.Called(words, multiples, limit)
	return args.Error(0)
}

func TestProcessMessage(t *testing.T) {
	repo := new(mockRepo)
	coreTest := core.NewCore[any](repo)

	tests := []struct {
		name      string
		words     []string
		multiples []int
		limit     int
		expected  string
		err       error
	}{
		{
			name:      "FizzBuzz case",
			multiples: []int{3, 5},
			words:     []string{"Fizz", "Buzz"},
			limit:     5,
			expected:  "[1 2 Fizz 4 Buzz]",
			err:       nil,
		},
		{
			name:      "Single number case",
			multiples: []int{2},
			words:     []string{"Even"},
			limit:     4,
			expected:  "[1 Even 3 Even]",
			err:       nil,
		},
		{
			name:      "Empty multiples and words",
			multiples: []int{},
			words:     []string{},
			limit:     3,
			expected:  "[1 2 3]",
			err:       nil,
		},
		{
			name:      "Limit zero",
			multiples: []int{3},
			words:     []string{"Fizz"},
			limit:     0,
			expected:  "[]",
			err:       nil,
		},
		{
			name:      "Multiple matches",
			multiples: []int{2, 4},
			words:     []string{"Two", "Four"},
			limit:     4,
			expected:  "[1 Two 3 TwoFour]",
			err:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("SaveMessage", tt.multiples, tt.words, tt.limit).Return(nil).Once()
			result, err := coreTest.ProcessMessage(tt.words, tt.multiples, tt.limit)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
