package main

import (
	"log"
	"time"

	"github.com/wallyqs/kit.go"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Starting...")

	// Create a NATS Component
	nc := kit.NewComponent()

	// TODO: Augmented behaviors for 'foo'

	// func(nc *nats.Conn, ctx context.Context, subj, reply string, data []byte)
	// nc.RegisterStreamWrapper("foo", func(data []byte){
	// 	log.Println("--- foo: ", data)
	// })

	err := nc.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// Bare NATS connection usage
	// nc.Subscribe("foo", func(m *nats.Msg){
	// 	log.Println("foo")
	// })
	// nc.Publish("foo", []byte("bar"))

	// Publishes on the foo subject, in case there is any middleware
	// for foo, it uses that to publish.

	// TODO: does
	// nc.Stream("foo").Publish([]byte("bar"))
	stream := nc.Stream("foo")

	// Need to configure the stream publisher generate once only.
	stream.Publish([]byte("bar"))

	// Returns a clone of the original event but with a different context.
	// stream.WithContext(ctx).Publish([]byte("bar"))

	// OnPublish each one of the callbacks get re executed.
	// Have to be able to pass a new context once again.
	svc := nc.Service("help")

	// Subscribe... and QueueSubscribe
	svc.Subscribe(func(m *nats.Msg){
		// TODO: But this does not apply the context back...
		m.Respond([]byte("pong"))
	})

	// Declare a reusable encoder.
	// svc := nc.Service("help").ApplyMiddleware(kit.OnSub(...), kit.OnPub(...))

	// calling it twice creates fails since already registered a subscription
	// to the service, can only have one type of subscription.
	// svc := nc.Service("help").Subscribe(func(m *nats.Msg){
	// })

	// svc := nc.Service("help").WithContext(ctx)

	resp, err := svc.Request([]byte("help please"), 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response:", resp)

	//
	// nc.Service("foo").WithContext(ctx).Request([]byte("bar"))
	select {}
}
