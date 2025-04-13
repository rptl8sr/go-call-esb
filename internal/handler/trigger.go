package handler

import (
	"encoding/json"
	"fmt"
)

type Source string

const (
	httpSource      Source = "http"
	timerSource     Source = "timer"
	unknownSource   Source = "unknown"
	notParsedSource Source = "not parsed"
)

// TimerEvent represents the structure of an event from a Yandex Cloud timer trigger.
type TimerEvent struct {
	Details struct {
		TriggerID string `json:"trigger_id"`
	} `json:"details"`
}

// HTTPEvent represents the structure of an event from a Yandex Cloud HTTP trigger.
type HTTPEvent struct {
	HTTPMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Url        string            `json:"url"`
}

// DetectTriggerType determines the type of trigger that invoked the function (timer or HTTP).
// Returns "timer", "http", or "unknown" if the event type is not recognized.
func DetectTriggerType(event interface{}) string {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return string(notParsedSource)
	}

	// TimerEvent
	var timerEvent TimerEvent
	err = json.Unmarshal(eventBytes, &timerEvent)
	if err == nil && timerEvent.Details.TriggerID != "" {
		return fmt.Sprintf("%s: %s", timerSource, timerEvent.Details.TriggerID)
	}

	// HTTPEvent
	var httpEvent HTTPEvent
	err = json.Unmarshal(eventBytes, &httpEvent)
	if err == nil && httpEvent.HTTPMethod != "" {
		return string(httpSource)
	}

	// Default
	return string(unknownSource)
}
