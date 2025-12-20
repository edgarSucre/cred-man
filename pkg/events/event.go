package events

import "context"

type Event interface {
	EventName() string
}

type Bus interface {
	Publish(context.Context, Event) error
}
