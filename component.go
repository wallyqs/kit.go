package kit

import (
	"time"

	"github.com/nats-io/nats.go"
)

// Component is oriented towards best practices of NATS usage
// and augmented behaviors on top of the NATS protocol.
// More or or less inspired in goals with the Clojure component
// library.
type Component struct {
	nc *nats.Conn

	// TODO: Add optional logger for example

	// TODO: There would be many type of firms here though.
	pubwrappers map[string]func(data []byte)
}

// Connect should have the same API as nats.Connect,
// options for the component should be done when
// creating the component.
func (c *Component) Connect(url string, options ...nats.Option) error {
	// TODO: Have some default callbacks implemented.

	nc, err := nats.Connect(url, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// TODO: Apply create all subscriptions, or services
	// then call flush.

	return nil
}

// TODO: NatsDefaultsOptions func

// In case NATS connection already present then this just starts.
// Assumes that already connected to NATS connected.
func (c *Component) Start() error {
	return nil
}

// These aren't worthy probably... should just use the ones
// from the library or use the services, stream abstractions.
// TODO: Probably only support async callbacks.
func (c *Component) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := c.nc.Subscribe(subj, cb)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// Publish uses the bare NATS connection to emit data.
func (c *Component) Publish(subj string, payload []byte) error {
	return c.nc.Publish(subj, payload)
}

// NOPE: don't do this, do allow using out of band the NATS conn.
// func (c *Component) NATS()

type Stream struct {
	subject string
	cbs     []func()
	nc      *nats.Conn
	sub     *nats.Subscription
}

func (e *Stream) Publish(payload []byte) error {
	// 
	// TODO: Apply the callbacks...
	// 
	return e.nc.Publish(e.subject, payload)
}

func (e *Stream) Subscribe(cb nats.MsgHandler) error {
	// 
	// TODO: Scaffolding the repo only but not what I want exactly...
	//
	sub, err := e.nc.Subscribe(e.subject, cb)
	if err != nil {
		return err
	}
	e.sub = sub

	return nil
}

// Stream returns an event that can be published.
func (c *Component) Stream(subj string) *Stream {
	// Pick the cbs that will be executed from this event.
	cbs := make([]func(), 0)
	return &Stream{
		subject: subj,
		cbs:     cbs,
		nc:      c.nc,
	}
}

type Service struct {
	subject string
	cbs     []func()
	nc      *nats.Conn
	sub     *nats.Subscription
}

func (e *Service) Subscribe(cb nats.MsgHandler) error {
	// 
	// TODO: Scaffolding the repo only but not what I want exactly...
	//
	sub, err := e.nc.Subscribe(e.subject, cb)
	if err != nil {
		return err
	}
	e.sub = sub

	return nil
}

func (e *Service) Request(payload []byte, timeout time.Duration) (*nats.Msg, error) {
	//
	// TODO: Apply the callbacks...
	//
	return e.nc.Request(e.subject, payload, timeout)
}

// Service returns an event that can be published.
func (c *Component) Service(subj string) *Service {
	// Pick the cbs that will be executed from this event.
	cbs := make([]func(), 0)
	return &Service{
		subject: subj,
		cbs:     cbs,
		nc:      c.nc,
	}
}

// NewComponent
//
// TODO: This only exists to avoid setting some private fields.
func NewComponent() *Component {
	return &Component{
		// TODO: but there can be a collection of wrappers though.
		pubwrappers: make(map[string]func([]byte)),
	}
}
