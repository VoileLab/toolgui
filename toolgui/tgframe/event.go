package tgframe

import (
	"encoding/json"
	"fmt"
)

// Event is the state change made by app user
type Event interface {
	ApplyState(state *State)
}

// EventType is the type of event
type EventType string

const (
	EventEmptyName EventType = ""
	EventClickName EventType = "click"
	EventInputName EventType = "input"
)

type EventStruct struct {
	Type EventType `json:"type"`
}

func ParseEvent(data []byte) (Event, error) {
	var event EventStruct
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case EventEmptyName:
		return &EventEmpty{}, nil
	case EventClickName:
		var eventClick EventClick
		err = json.Unmarshal(data, &eventClick)
		if err != nil {
			return nil, err
		}
		return &eventClick, nil
	case EventInputName:
		var eventInput EventInput
		err = json.Unmarshal(data, &eventInput)
		if err != nil {
			return nil, err
		}
		return &eventInput, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// EventEmpty is the event of an empty event, rerun button will send this
type EventEmpty struct {
}

func (e *EventEmpty) ApplyState(*State) {
}

// EventClick is the event of a button click event
type EventClick struct {
	ID string `json:"id"`
}

func (e *EventClick) ApplyState(state *State) {
	state.SetClickID(e.ID)
}

// EventInput is the event of a input event
// it's used for all input component
type EventInput struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

func (e *EventInput) ApplyState(state *State) {
	state.Set(e.ID, e.Value)
}
