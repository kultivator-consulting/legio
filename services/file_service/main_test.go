package main

import (
	"log"
	"os"
	"testing"
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func TestLoadEnv(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	table := []struct {
		name     string
		input    string
		expected string
	}{
		{"APP_ENV_unset", "--unset--", ".env"},
		{"APP_ENV_empty_env", "", ".env"},
		{"APP_ENV_development", "development", ".env.development"},
		{"APP_ENV_test", "test", ".env.test"},
		{"APP_ENV_staging", "staging", ".env.staging"},
		{"APP_ENV_production", "production", ".env.production"},
		{"APP_ENV_garbage", "garbage", ".env.garbage"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			teardownTest := setupTest(t)
			defer teardownTest(t)

			filename := ""

			if tc.input == "--unset--" {
				err := os.Unsetenv("APP_ENV")
				if err != nil {
					t.Errorf("LoadEnv() error = %v", err)
				}
			} else {
				t.Setenv("APP_ENV", tc.input)
			}

			spyGodotenvLoad = func(envFilenames ...string) error {
				if len(envFilenames) > 0 {
					filename = envFilenames[0]
				}
				return nil
			}

			err := LoadEnv()
			if err != nil {
				t.Errorf("LoadEnv() error = %v", err)
			}
			if filename != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, filename)
			}
		})
	}
}
