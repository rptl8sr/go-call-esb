package handler

import (
	"testing"
)

func TestDetectTriggerType(t *testing.T) {
	tests := []struct {
		name     string
		event    interface{}
		expected string
	}{
		{
			name: "Valid Timer Event",
			event: TimerEvent{
				Details: struct {
					TriggerID string `json:"trigger_id"`
				}{TriggerID: "123"},
			},
			expected: "timer: 123",
		},
		{
			name: "Valid HTTP Event",
			event: HTTPEvent{
				HTTPMethod: "GET",
				Url:        "/api/v1/resource",
			},
			expected: "http",
		},
		{
			name: "Unknown Event Type",
			event: struct {
				SomeField string
			}{
				SomeField: "random_value",
			},
			expected: "unknown",
		},
		{
			name:     "Empty Struct",
			event:    struct{}{},
			expected: "unknown",
		},
		{
			name: "Invalid Timer Event",
			event: TimerEvent{
				Details: struct {
					TriggerID string `json:"trigger_id"`
				}{TriggerID: ""},
			},
			expected: "unknown",
		},
		{
			name: "Invalid HTTP Event (Empty HTTPMethod)",
			event: HTTPEvent{
				HTTPMethod: "",
				Url:        "/api/v1/resource",
			},
			expected: "unknown",
		},
		{
			name:     "Serialization Error",
			event:    make(chan int),
			expected: "not parsed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTriggerType(tt.event)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
