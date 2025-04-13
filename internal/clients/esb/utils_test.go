package esb

import (
	"reflect"
	"testing"
)

func Test_buildHeaders(t *testing.T) {
	tests := []struct {
		name           string
		expected       map[string]string
		expectedLength int
	}{
		{
			name:           "Get headers",
			expected:       map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
			expectedLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHeaders()
			if len(result) != tt.expectedLength {
				t.Errorf("buildHeaders(%d), expected %d", len(result), tt.expectedLength)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("buildHeaders(%v), expected %v", result, tt.expected)
			}
		})
	}
}
