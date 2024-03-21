// package config

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestErrorLoadConfig(t *testing.T) {
// 	const testConfigPath = "../salah"

// 	err := LoadConfig(testConfigPath)

// 	assert.Error(t, err)

// 	assert.Empty(t, Configuration)

// }

// func TestLoadConfig(t *testing.T) {
// 	const testConfigPath = "../"

// 	err := LoadConfig(testConfigPath)

// 	assert.NoError(t, err, "Tidak seharusnya terjadi error saat memuat konfigurasi")

// 	expectedHost := "localhost"
// 	actualHost := Configuration.Mysql.Host
// 	assert.Equal(t, expectedHost, actualHost, "Should equal")

// }

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoading(t *testing.T) {
	testCases := []struct {
		name             string
		configPath       string
		expectedError    bool
		expectedHost     string
		expectedEmptyMap bool
	}{
		{
			name:             "Error loading config",
			configPath:       "../salah",
			expectedError:    true,
			expectedHost:     "",
			expectedEmptyMap: true,
		},
		{
			name:             "Successful config loading",
			configPath:       "../",
			expectedError:    false,
			expectedHost:     "localhost",
			expectedEmptyMap: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := LoadConfig(tc.configPath)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHost, Configuration.Mysql.Host)
			}

			if tc.expectedEmptyMap {
				assert.Empty(t, Configuration)
			}
		})
	}
}
