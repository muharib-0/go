package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      string
		expected int
	}{
		{
			name:     "Born today",
			dob:      time.Now().Format("2006-01-02"),
			expected: 0,
		},
		{
			name:     "Born 10 years ago exact",
			dob:      time.Now().AddDate(-10, 0, 0).Format("2006-01-02"),
			expected: 10,
		},
		{
			name:     "Born 20 years and 1 day ago",
			dob:      time.Now().AddDate(-20, 0, -1).Format("2006-01-02"),
			expected: 20,
		},
		{
			name:     "Born 20 years minus 1 day ago (birthday tomorrow)",
			dob:      time.Now().AddDate(-20, 0, 1).Format("2006-01-02"),
			expected: 19,
		},
		{
			name:     "Future date",
			dob:      time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dob, err := time.Parse("2006-01-02", tt.dob)
			assert.NoError(t, err)

			age := CalculateAge(dob)
			assert.Equal(t, tt.expected, age)
		})
	}
}
