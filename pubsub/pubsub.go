package pubsub

import (
	"context"

	"github.com/lovoo/goka"
)

// EventHandler is an event handler. It will be called after message is arrived to consumer
type EventHandler interface {
	Handle(ctx context.Context, message interface{}) (err error)
}

//  GokaEventHandler is an event handler for goka context. It will be called after message is arrived to consumer
type GokaEventHandler interface {
	Handle(ctx goka.Context, message interface{})
}

// Publisher is a collection of behavior of a publisher
type Publisher interface {
	// Will send the message to the assigned topic.
	Send(ctx context.Context, key string, message interface{}) (err error)
	Close() (err error)
}

// Subscriber is a collection of behavior of a subscriber
type Subscriber interface {
	Subscribe()
	Close() (err error)
}

// ViewTable is a collection of behavior of a view table
type ViewTable interface {
	Open()
	Get(key string) (data interface{}, err error)
	Close()
}

// GokaCodec is a collection of behavior of goka codec
type GokaCodec interface {
	Encode(value interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}

type GokaContext interface {
	goka.Context
}
