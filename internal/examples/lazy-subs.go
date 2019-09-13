package main

import (
	"log"
	"time"

	"github.com/wallyqs/kit.go"
	"github.com/nats-io/nats.go"
)

func main() {
	nc := kit.NewComponent()

	// Creates a stream object internally.
	nc.Stream("foo").Subscribe(func(m *nats.Msg){
		log.Println("->>", m)
	})

	// Creates a service object internally.
	// These have to be an interface to allow multiple type of usage.
	nc.Service("help").Subscribe(func(m *nats.Msg){
		log.Println("-->", m)

		// Respond skips the filters so can't use the object.
		m.Respond([]byte("pong"))
	})

	// Sends all subscriptions and does a flush.
	// If there is any error during connect, like an auth error
	// on a subscription then 
	err := nc.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// 
	// Can't create new services after being connected? Why not?
	// HTTP Servers can't create new handlers after starting,
	// can create new subscriptions via subs though, so not really
	// constrained either.
	// 

	// Can reuse the middleware that has been installed so far.
	svc := nc.Service("help")
	resp, err := svc.Request([]byte("help please"), 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response:", resp)

	//
	// nc.Service("foo").WithContext(ctx).Request([]byte("bar"))
	select {}
}
