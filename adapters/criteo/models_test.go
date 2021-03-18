package criteo

import (
	"testing"
)

func TestGetDeviceType(t *testing.T) {

	// Setup:
	deviceTypeCases := []struct {
		deviceType string
		expected   string
	}{
		{"ios", "idfa"},
		{"Ios", "idfa"},
		{"IOS", "idfa"},
		{"android", "gaid"},
		{"unknown", "unknown"},
		{"", "unknown"},
		{"qwerty", "unknown"},
		{"qWerty", "unknown"},
		{"abc", "unknown"},
	}

	for _, uc := range deviceTypeCases {

		// Execute:
		result := getDeviceType(uc.deviceType)

		// Verify:
		if uc.expected != result {
			t.Errorf("Bad getDeviceType for '%s'. Expected: %s, got %s", uc.deviceType, uc.expected, result)
		}
	}
}
