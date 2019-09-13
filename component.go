package kit

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Component is oriented towards best practices of NATS usage
// and augmented behaviors on top of the NATS protocol.
// More or or less inspired in goals with the Clojure component
// library.
type Component struct {
	nc *nats.Conn

	// Note: Subscriptions are singletons.
	streams  map[string]*Stream
	services map[string]*Service

	// connected
	connected bool
}

// Connect should have the same API as nats.Connect,
// options for the component should be done when
// creating the component.
func (c *Component) Connect(url string, options ...nats.Option) error {
	// TODO: Does require the error callback to be able to handle auth errors,
	// and issues with permissions.  Although internally the server does disconnect
	// the client from subscribing I think...

	nc, err := nats.Connect(url, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// TODO: Apply create all subscriptions, or services
	// then call flush.
	for subject, stream := range c.streams {
		log.Println("subscribing", stream, stream.cb)
		sub, err := nc.Subscribe(subject, stream.cb)
		if err != nil {
			log.Println(err)
		}
		stream.sub = sub
	}

	for subject, svc := range c.services {
		log.Println("subscribing", svc, svc.cb)
		sub, err := nc.Subscribe(subject, svc.cb)
		if err != nil {
			log.Println(err)
		}
		svc.sub = sub
	}

	return nil
}

// In case NATS connection already present then this just starts.
// Assumes that already connected to NATS connected.
func (c *Component) Start() error {
	return nil
}

// These aren't worthy probably... should just use the ones
// from the library or use the services, stream abstractions.
// TODO: Probably only support async callbacks for now.
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

// Service...
type Service struct {
	subject string
	hooks   []func()
	nc      *nats.Conn
	sub     *nats.Subscription
	comp    *Component
	cb      nats.MsgHandler
}

// Subscribe...
func (svc *Service) Subscribe(cb nats.MsgHandler) error {
	// Original callback
	svc.cb = cb
	if svc.comp.nc != nil {
		//
		// TODO: Scaffolding the repo only but not what I want exactly...
		//
		sub, err := svc.nc.Subscribe(svc.subject, cb)
		if err != nil {
			return err
		}
		svc.sub = sub
	}

	// Register the service...
	svc.comp.services[svc.subject] = svc

	return nil
}

// Request...
func (svc *Service) Request(payload []byte, timeout time.Duration) (*nats.Msg, error) {
	//
	// TODO: Apply the callbacks...
	//
	return svc.nc.Request(svc.subject, payload, timeout)
}

// Service returns an event that can be published...
func (c *Component) Service(subj string) *Service {
	// Pick the cbs that will be executed from this event.
	return &Service{
		subject: subj,
		hooks:   make([]func(), 0),
		nc:      c.nc,
		comp:    c,
	}
}

// NewComponent...
func NewComponent() *Component {
	return &Component{
		streams:  make(map[string]*Stream),
		services: make(map[string]*Service),
	}
}
