package controller_test

import (
	"context"
	"encoding/json"
	"fizzbuzz/internal/controller"
	"fizzbuzz/internal/types"
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

func (m *MockCore) GetStatsWords() ([]types.StatsWordsResult, error) {
	args := m.Called()
	return args.Get(0).([]types.StatsWordsResult), args.Error(1)
}

func (m *MockCore) GetStatsParameters() ([]types.StatsParameters, error) {
	args := m.Called()
	return args.Get(0).([]types.StatsParameters), args.Error(1)
}

func (m *MockCore) GetTotalRequests() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func TestFizzBuzzController(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	mockCore := new(MockCore)
	testController := controller.NewFizzBuzzController(mockCore)

	tests := []struct {
		name           string
		input          controller.FizzBuzzRequest
		queryParams    string
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
			queryParams: "multiples=3,5&words=Fizz,Buzz&limit=5",
			mockSetup: func() {
				mockCore.On("ProcessMessage", []string{"Fizz", "Buzz"}, []int{3, 5}, 5).
					Return("[1 2 Fizz 4 Buzz]", nil).Once()
			},
			expectedStatus: 200,
			expectedBody:   gin.H{"result": "[1 2 Fizz 4 Buzz]"},
		},
		{
			name: "Invalid JSON",
			input: controller.FizzBuzzRequest{
				Multiples: []int{3, 5},
				Words:     []string{"Fizz", "Buzz"},
				Limit:     0, // Invalid limit
			},
			queryParams:    "multiples=3,5&words=Fizz,Buzz&limit=0",
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
			queryParams:    "multiples=3,5&words=Fizz&limit=5",
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
			queryParams: "multiples=3,5&words=Fizz,Buzz&limit=5",
			mockSetup: func() {
				mockCore.On("ProcessMessage", []string{"Fizz", "Buzz"}, []int{3, 5}, 5).
					Return("", assert.AnError).Once()
			},
			expectedStatus: 500,
			expectedBody:   gin.H{"error": assert.AnError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/fizzbuzz", testController.FizzBuzz)

			tt.mockSetup()

			w := httptest.NewRecorder()
			url := "/fizzbuzz?" + tt.queryParams
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var actualBody gin.H
			json.Unmarshal(w.Body.Bytes(), &actualBody)
			assert.Equal(t, tt.expectedBody, actualBody)

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
				mockCore.On("GetStatsParameters").Return(
					[]types.StatsParameters{
						{Words: []string{"Fizz", "Buzz"}, Multiples: []int{3, 5}, Limit: 15, TotalRequests: 5},
					}, nil).Once()
				mockCore.On("GetTotalRequests").Return(10, nil).Once()
			},
			expectedStatus: 200,
			expectedBody: controller.StatResp{
				TotalRequests: 10,
				RequestStats: []types.StatsParameters{
					{Words: []string{"Fizz", "Buzz"}, Multiples: []int{3, 5}, Limit: 15, TotalRequests: 5},
				},
			},
		},
		{
			name: "Core error on stats",
			mockSetup: func() {
				mockCore.On("GetStatsParameters").Return([]types.StatsParameters{}, assert.AnError).Once()
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
