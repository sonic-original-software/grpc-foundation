package server

import (
	"testing"
	"time"
)

func TestAddress(t *testing.T) {
	t.Run("returns default when environment variable not set", func(t *testing.T) {
		result := Address()
		if result != DefaultAddress {
			t.Errorf("expected %q, got %q", DefaultAddress, result)
		}
	})

	t.Run("returns environment variable when set", func(t *testing.T) {
		t.Setenv(EnvServerAddress, ":9090")
		result := Address()
		if result != ":9090" {
			t.Errorf("expected %q, got %q", ":9090", result)
		}
	})
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("returns environment variable when set", func(t *testing.T) {
		t.Setenv("TEST_VAR", "custom-value")
		result := getEnvOrDefault("TEST_VAR", "default-value")
		if result != "custom-value" {
			t.Errorf("expected %q, got %q", "custom-value", result)
		}
	})

	t.Run("returns default when environment variable not set", func(t *testing.T) {
		result := getEnvOrDefault("NONEXISTENT_VAR", "default-value")
		if result != "default-value" {
			t.Errorf("expected %q, got %q", "default-value", result)
		}
	})

	t.Run("returns default when environment variable is empty", func(t *testing.T) {
		t.Setenv("EMPTY_VAR", "")
		result := getEnvOrDefault("EMPTY_VAR", "default-value")
		if result != "default-value" {
			t.Errorf("expected %q, got %q", "default-value", result)
		}
	})
}

func TestGetEnvDurationOrDefault(t *testing.T) {
	t.Run("returns parsed duration when valid", func(t *testing.T) {
		t.Setenv("DURATION_VAR", "10m")
		result := getEnvDurationOrDefault("DURATION_VAR", 5*time.Minute)
		if result != 10*time.Minute {
			t.Errorf("expected %v, got %v", 10*time.Minute, result)
		}
	})

	t.Run("returns default when not set", func(t *testing.T) {
		result := getEnvDurationOrDefault("NONEXISTENT_DURATION", 5*time.Minute)
		if result != 5*time.Minute {
			t.Errorf("expected %v, got %v", 5*time.Minute, result)
		}
	})

	t.Run("returns default when invalid format", func(t *testing.T) {
		t.Setenv("INVALID_DURATION", "not-a-duration")
		result := getEnvDurationOrDefault("INVALID_DURATION", 5*time.Minute)
		if result != 5*time.Minute {
			t.Errorf("expected %v, got %v", 5*time.Minute, result)
		}
	})

	t.Run("handles various duration formats", func(t *testing.T) {
		t.Setenv("SECONDS_VAR", "30s")
		result := getEnvDurationOrDefault("SECONDS_VAR", time.Minute)
		if result != 30*time.Second {
			t.Errorf("expected %v, got %v", 30*time.Second, result)
		}
	})
}

func TestGetEnvIntOrDefault(t *testing.T) {
	t.Run("returns parsed int when valid", func(t *testing.T) {
		t.Setenv("INT_VAR", "42")
		result := getEnvIntOrDefault("INT_VAR", 10)
		if result != 42 {
			t.Errorf("expected %d, got %d", 42, result)
		}
	})

	t.Run("returns default when not set", func(t *testing.T) {
		result := getEnvIntOrDefault("NONEXISTENT_INT", 10)
		if result != 10 {
			t.Errorf("expected %d, got %d", 10, result)
		}
	})

	t.Run("returns default when invalid format", func(t *testing.T) {
		t.Setenv("INVALID_INT", "not-an-int")
		result := getEnvIntOrDefault("INVALID_INT", 10)
		if result != 10 {
			t.Errorf("expected %d, got %d", 10, result)
		}
	})

	t.Run("handles negative numbers", func(t *testing.T) {
		t.Setenv("NEGATIVE_INT", "-5")
		result := getEnvIntOrDefault("NEGATIVE_INT", 10)
		if result != -5 {
			t.Errorf("expected %d, got %d", -5, result)
		}
	})
}
