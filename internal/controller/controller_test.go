package controller_test

import (
	"bytes"
	"encoding/json"
	"fizzbuzz/internal/controller"
	"fizzbuzz/internal/types"
	"fmt"
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

func (m *MockCore) GetStatsWords() ([]types.StatsByKeyResult, error) {
	args := m.Called()
	return args.Get(0).([]types.StatsByKeyResult), args.Error(1)
}

func (m *MockCore) GetStatsParameters() (types.StatsParameters, error) {
	args := m.Called()
	return args.Get(0).(types.StatsParameters), args.Error(1)
}

func (m *MockCore) GetTotalRequests() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
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
				mockCore.On("ProcessMessage", mock.Anything, mock.Anything, mock.Anything).
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
			params := "multiples=3,5&words=Fizz,Buzz&limit=5"
			c.Request, _ = http.NewRequest("POST", fmt.Sprintf("/fizzbuzz?%s", params), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("Accept", "application/json")

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

func TestStats(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new mock Core
	mockCore := new(MockCore)
	testController := controller.NewFizzBuzzController(mockCore)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Valid stats request",
			mockSetup: func() {
				mockCore.On("GetStats").Return([]types.StatsByKeyResult{
					{Key: "Fizz", Count: 10},
					{Key: "Buzz", Count: 5},
				}, nil).Once()
				mockCore.On("GetTotalRequests").Return(10, nil).Once()
			},
			expectedStatus: 200,
			expectedBody: controller.StatResp{
				TotalRequests: 10,
				Stats: []types.StatsByKeyResult{
					{Key: "Fizz", Count: 10},
					{Key: "Buzz", Count: 5},
				},
			},
		},
		{
			name: "Core error on stats",
			mockSetup: func() {
				mockCore.On("GetStats").Return(nil, assert.AnError).Once()
				mockCore.On("GetTotalRequests").Return(0, nil).Once()
			},
			expectedStatus: 500,
			expectedBody:   gin.H{"error": assert.AnError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/stats", nil)
			c.Request.Header.Set("Content-Type", "application/json")

			testController.Stats(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == 200 {
				var result controller.StatResp
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, tt.expectedBody, result)
			}
		})
	}
}
