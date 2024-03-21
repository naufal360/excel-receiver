package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateReqID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "Generate request ID with default length (6)",
			expectedLength: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqID := GenerateReqID()

			assert.Equal(t, tc.expectedLength, len(reqID), "Generated request ID length does not match the expected length")
		})
	}
}
