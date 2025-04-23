package domain_test

import (
	"beetle/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationQuery_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		query    domain.PaginationQuery
		expected int
	}{
		{
			name: "first page",
			query: domain.PaginationQuery{
				Limit: 10,
				Page:  1,
			},
			expected: 0,
		},
		{
			name: "second page",
			query: domain.PaginationQuery{
				Limit: 10,
				Page:  2,
			},
			expected: 10,
		},
		{
			name: "third page with custom limit",
			query: domain.PaginationQuery{
				Limit: 25,
				Page:  3,
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := tt.query.GetOffset()
			assert.Equal(t, tt.expected, offset)
		})
	}
}

func TestPaginationQuery_CreateResults(t *testing.T) {
	tests := []struct {
		name           string
		query          domain.PaginationQuery
		expectedOffset int
	}{
		{
			name: "create results with default values",
			query: domain.PaginationQuery{
				Limit: 10,
				Page:  1,
			},
			expectedOffset: 0,
		},
		{
			name: "create results with custom values",
			query: domain.PaginationQuery{
				Limit: 20,
				Page:  2,
			},
			expectedOffset: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, offset := tt.query.CreateResults()

			assert.Equal(t, tt.query, results.PaginationQuery)
			assert.Equal(t, tt.expectedOffset, offset)
			assert.Nil(t, results.Data)
			assert.Equal(t, 0, results.Total)
		})
	}
}

func TestPaginationConstants(t *testing.T) {
	assert.Equal(t, domain.PageLimitMax, 120)
	assert.Equal(t, domain.PageLimitMin, 1)
	assert.Equal(t, domain.PageLimitDefault, 12)
}
