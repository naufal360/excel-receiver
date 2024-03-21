package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken_IsExpired(t *testing.T) {
	testCases := []struct {
		name        string
		token       Token
		expectation bool
	}{
		{
			name: "Valid token",
			token: Token{
				ExpiredAt: time.Now().Add(1 * time.Hour),
			},
			expectation: true, // Expecting this token to be expired
		},
		{
			name: "Expired token",
			token: Token{
				ExpiredAt: time.Now().Add(10 * time.Hour),
			},
			expectation: false, // Expecting this token not to be expired
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectation, tc.token.IsExpired())
		})
	}
}
