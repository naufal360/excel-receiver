package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	testCases := []struct {
		name       string
		paths      []string
		createErr  bool
		expectErr  bool
		expectDirs bool
	}{
		{
			name:       "Create single directory",
			paths:      []string{"testdir"},
			createErr:  false,
			expectErr:  false,
			expectDirs: true,
		},
		{
			name:       "Create multiple directories",
			paths:      []string{"testdir1", "testdir2"},
			createErr:  false,
			expectErr:  false,
			expectDirs: true,
		},
		{
			name:       "Create existing directory",
			paths:      []string{"testdir"},
			createErr:  false,
			expectErr:  false,
			expectDirs: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up any previously created directories
			defer os.RemoveAll("testdir")
			defer os.RemoveAll("testdir1")
			defer os.RemoveAll("testdir2")

			// Create directories
			for _, path := range tc.paths {
				err := CreateDirectory(path)
				if tc.createErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}

			// Check if directories exist
			for _, path := range tc.paths {
				_, err := os.Stat(path)
				if tc.expectDirs {
					assert.NoError(t, err)
				} else {
					assert.True(t, os.IsNotExist(err))
				}
			}
		})
	}
}
