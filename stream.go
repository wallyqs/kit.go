package kit

import (
	"github.com/nats-io/nats.go"
)

// Stream...
type Stream struct {
	subject string
	hooks   []func()
	nc      *nats.Conn
	sub     *nats.Subscription
	comp    *Component
	cb      nats.MsgHandler
}

// Publish...
func (stream *Stream) Publish(payload []byte) error {
	//
	// TODO: Apply the callbacks...
	//
	return stream.nc.Publish(stream.subject, payload)
}

// Subscribe...
func (stream *Stream) Subscribe(cb nats.MsgHandler) error {
	stream.cb = cb
	if stream.comp.nc != nil {
		// Should be deferred then
		//
		// TODO: Scaffolding the repo only but not what I want exactly...
		//
		sub, err := stream.nc.Subscribe(stream.subject, cb)

		//
		// TODO: errors on subscribe are fatal? should be disconnected.
		//
		if err != nil {
			return err
		}
		stream.sub = sub
	}

	// Register the stream.
	stream.comp.streams[stream.subject] = stream

	return nil
}

// Stream returns an event that can be published.
func (c *Component) Stream(subj string) *Stream {
	// Pick the cbs that will be executed from this event.
	return &Stream{
		subject: subj,
		hooks:   make([]func(), 0),
		nc:      c.nc,
		comp:    c,
	}
}
