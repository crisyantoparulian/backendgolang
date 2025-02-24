package config

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		cleanupEnv     func()
		expectedConfig *Config
		expectedError  error
	}{
		{
			name: "Success - Load Config from Environment Variables",
			setupEnv: func() {
				os.Setenv("APP_PORT", "8080")
				os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname")
			},
			cleanupEnv: func() {
				os.Unsetenv("APP_PORT")
				os.Unsetenv("DATABASE_URL")
			},
			expectedConfig: &Config{
				App: App{
					Port: 8080,
				},
				Database: Database{
					PostgreDSN: "postgres://user:password@localhost:5432/dbname",
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid APP_PORT",
			setupEnv: func() {
				os.Setenv("APP_PORT", "invalid")
				os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname")
			},
			cleanupEnv: func() {
				os.Unsetenv("APP_PORT")
				os.Unsetenv("DATABASE_URL")
			},
			expectedConfig: nil,
			expectedError:  fmt.Errorf("invalid APP_PORT in .env: strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
		{
			name: "Success - Load Config Only Once",
			setupEnv: func() {
				os.Setenv("APP_PORT", "8080")
				os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname")
			},
			cleanupEnv: func() {
				os.Unsetenv("APP_PORT")
				os.Unsetenv("DATABASE_URL")
			},
			expectedConfig: &Config{
				App: App{
					Port: 8080,
				},
				Database: Database{
					PostgreDSN: "postgres://user:password@localhost:5432/dbname",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			tt.setupEnv()

			// Reset the sync.Once to ensure LoadConfig runs again
			once = sync.Once{}
			cfg = nil

			// Call LoadConfig
			if tt.expectedError != nil {
				assert.PanicsWithError(t, tt.expectedError.Error(), func() {
					LoadConfig()
				})
			} else {
				config := LoadConfig()
				assert.Equal(t, tt.expectedConfig, config, "loaded config should match expected config")
			}

			// Cleanup environment variables
			tt.cleanupEnv()
		})
	}
}

func TestLoadConfig_LoadOnlyOnce(t *testing.T) {
	// Setup environment variables
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname")
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DATABASE_URL")
	}()

	// Reset the sync.Once to ensure LoadConfig runs again
	once = sync.Once{}
	cfg = nil

	// Call LoadConfig twice
	config1 := LoadConfig()
	config2 := LoadConfig()

	// Ensure both calls return the same instance
	assert.Equal(t, config1, config2, "LoadConfig should return the same instance on subsequent calls")
}
