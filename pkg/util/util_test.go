package util

import "testing"

func TestSetLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "negative limit",
			input:    -5,
			expected: 10,
		},
		{
			name:     "zero limit",
			input:    0,
			expected: 10,
		},
		{
			name:     "valid limit",
			input:    25,
			expected: 25,
		},
		{
			name:     "limit at upper bound",
			input:    100,
			expected: 100,
		},
		{
			name:     "limit exceeds upper bound",
			input:    150,
			expected: 100,
		},
		{
			name:     "limit of 1",
			input:    1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetLimit(tt.input)
			if got != tt.expected {
				t.Errorf("SetLimit(%d) = %d, expected %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		page     int
		expected int
	}{
		{
			name:     "page 1 with limit 10",
			limit:    10,
			page:     1,
			expected: 0,
		},
		{
			name:     "page 2 with limit 10",
			limit:    10,
			page:     2,
			expected: 10,
		},
		{
			name:     "page 3 with limit 20",
			limit:    20,
			page:     3,
			expected: 40,
		},
		{
			name:     "zero page (should default to 1)",
			limit:    10,
			page:     0,
			expected: 0,
		},
		{
			name:     "negative page (should default to 1)",
			limit:    10,
			page:     -1,
			expected: 0,
		},
		{
			name:     "limit exceeds 100 (should cap at 100)",
			limit:    200,
			page:     2,
			expected: 100,
		},
		{
			name:     "limit below 1 (should default to 10)",
			limit:    -5,
			page:     3,
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateOffset(tt.limit, tt.page)
			if got != tt.expected {
				t.Errorf("CalculateOffset(%d, %d) = %d, expected %d", tt.limit, tt.page, got, tt.expected)
			}
		})
	}
}
