package controller_test

import (
	"bytes"
	"encoding/json"
	"fizzbuzz/internal/controller"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCore is a mock implementation of the Core interface
type MockCore struct {
	mock.Mock
}

func (m *MockCore) ProcessMessage(words []string, values []int, limit int) (string, error) {
	args := m.Called(words, values, limit)
	return args.String(0), args.Error(1)
}

func TestFizzBuzzController(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new mock Core
	mockCore := new(MockCore)
	testController := controller.NewFizzBuzzController(mockCore)

	tests := []struct {
		name           string
		input          controller.FizzBuzzRequest
		mockSetup      func()
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Valid request",
			input: controller.FizzBuzzRequest{
				Multiples: []int{3, 5},
				Words:     []string{"Fizz", "Buzz"},
				Limit:     5,
			},
			mockSetup: func() {
				mockCore.On("ProcessMessage", []int{3, 5}, []string{"Fizz", "Buzz"}, 5).
					Return("[1 2 Fizz 4 Buzz]", nil).Once()
			},
			expectedStatus: 200,
			expectedBody:   "[1 2 Fizz 4 Buzz]",
		},
		{
			name: "Invalid JSON",
			input: controller.FizzBuzzRequest{
				Multiples: []int{3, 5},
				Words:     []string{"Fizz", "Buzz"},
				Limit:     0, // Invalid limit
			},
			mockSetup:      func() {},
			expectedStatus: 400,
			expectedBody:   gin.H{"error": "params multiples, words and limit are required and must be valid"},
		},
		{
			name: "Mismatched array lengths",
			input: controller.FizzBuzzRequest{
				Multiples: []int{3, 5},
				Words:     []string{"Fizz"},
				Limit:     5,
			},
			mockSetup:      func() {},
			expectedStatus: 400,
			expectedBody:   gin.H{"error": "multiples and words arrays must have the same length"},
		},
		{
			name: "Core error",
			input: controller.FizzBuzzRequest{
				Multiples: []int{3, 5},
				Words:     []string{"Fizz", "Buzz"},
				Limit:     5,
			},
			mockSetup: func() {
				mockCore.On("ProcessMessage", []int{3, 5}, []string{"Fizz", "Buzz"}, 5).
					Return("", assert.AnError).Once()
			},
			expectedStatus: 500,
			expectedBody:   gin.H{"error": assert.AnError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup()

			// Create a new gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create JSON body
			body, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest("POST", "/fizzbuzz", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the controller
			testController.FizzBuzz(c)

			// Assert results
			assert.Equal(t, tt.expectedStatus, w.Code)
			var actualBody any
			if tt.expectedStatus == 200 {
				var result string
				json.Unmarshal(w.Body.Bytes(), &result)
				actualBody = result
			} else {
				var result gin.H
				json.Unmarshal(w.Body.Bytes(), &result)
				actualBody = result
			}
			assert.Equal(t, tt.expectedBody, actualBody)

			// Verify mock expectations
			mockCore.AssertExpectations(t)
		})
	}
}
