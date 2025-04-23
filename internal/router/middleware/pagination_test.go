package middleware_test

import (
	"beetle/internal/domain"
	"beetle/internal/handler"
	"beetle/internal/router/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name          string
		queryParams   string
		expectedPage  int
		expectedLimit int
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "default values",
			queryParams:   "",
			expectedPage:  1,
			expectedLimit: domain.PageLimitDefault,
			expectedError: false,
		},
		{
			name:          "valid page and limit",
			queryParams:   "?page=2&limit=20",
			expectedPage:  2,
			expectedLimit: 20,
			expectedError: false,
		},
		{
			name:          "limit too high",
			queryParams:   "?page=1&limit=121",
			expectedPage:  0,
			expectedLimit: 0,
			expectedError: true,
			errorMessage:  "limit must be between 1 and 120",
		},
		{
			name:          "limit too low",
			queryParams:   "?page=1&limit=0",
			expectedPage:  0,
			expectedLimit: 0,
			expectedError: true,
			errorMessage:  "limit must be between 1 and 120",
		},
		{
			name:          "page too low",
			queryParams:   "?page=0&limit=10",
			expectedPage:  0,
			expectedLimit: 0,
			expectedError: true,
			errorMessage:  "page must be greater than 0",
		},
		{
			name:          "invalid page format",
			queryParams:   "?page=abc&limit=10",
			expectedPage:  1,
			expectedLimit: 10,
			expectedError: false,
		},
		{
			name:          "invalid limit format",
			queryParams:   "?page=1&limit=xyz",
			expectedPage:  1,
			expectedLimit: domain.PageLimitDefault,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			hc := &handler.Context{Context: c}

			// Test
			err := middleware.Paginate(hc)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, he.Code)
				assert.Equal(t, tt.errorMessage, he.Message)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, hc.PaginationQuery)
				assert.Equal(t, tt.expectedPage, hc.PaginationQuery.Page)
				assert.Equal(t, tt.expectedLimit, hc.PaginationQuery.Limit)
			}
		})
	}
}

func TestParseIntQueryParam(t *testing.T) {
	tests := []struct {
		name           string
		paramName      string
		paramValue     string
		fallback       int
		expectedResult int
	}{
		{
			name:           "valid integer",
			paramName:      "test",
			paramValue:     "42",
			fallback:       10,
			expectedResult: 42,
		},
		{
			name:           "empty value",
			paramName:      "test",
			paramValue:     "",
			fallback:       10,
			expectedResult: 10,
		},
		{
			name:           "invalid format",
			paramName:      "test",
			paramValue:     "abc",
			fallback:       10,
			expectedResult: 10,
		},
		{
			name:           "negative number",
			paramName:      "test",
			paramValue:     "-5",
			fallback:       10,
			expectedResult: -5,
		},
		{
			name:           "zero",
			paramName:      "test",
			paramValue:     "0",
			fallback:       10,
			expectedResult: 0,
		},
		{
			name:           "large number",
			paramName:      "test",
			paramValue:     "999999",
			fallback:       10,
			expectedResult: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.paramName+"="+tt.paramValue, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			hc := &handler.Context{Context: c}

			// Test
			result := middleware.ParseIntQueryParam(hc, tt.paramName, tt.fallback)

			// Assertions
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
