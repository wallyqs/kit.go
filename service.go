package kit

import (
	"time"

	"github.com/nats-io/nats.go"
)

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
