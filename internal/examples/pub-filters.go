package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/wallyqs/kit.go"
)

func main() {
	// Disable echo by default?
	nc := kit.NewComponent()

	// Sets globally some type of filter on any subject that matches
	// the wildcard?  Too much overhead.
	// nc.RegisterPublishFilter()

	// Creates a stream object internally.
	stream := nc.Stream("foo")

	// Modifies any message published.  Applied per message.
	// TODO: Should be able to create default publish filters for topics

	// 
	// stream.PublishFilter(ctx context.Context, func (data []byte) []byte {
	//
	// if txn := newrelic.FromContext(ctx); nil != txn {
	// 	seg := newrelic.ExternalSegment{
	// 		StartTime: newrelic.StartSegmentNow(txn),
	// 		URL:       nc.ConnectedUrl(),
	// 		Procedure: "Publish/" + subj,
	// 		Library:   "NATS",
	// 	}   
	// 	defer seg.End()
	// }   
	//

	// Then middleware take a stream or request, and apply the filters.
	// Encoding middleware takes the stream and does encoding decoding.
	// jsonencoder.ApplyFilters(stream)
	
	stream.PublishFilter(func (data []byte) []byte {
		log.Println("::::: Filter ::::", string(data))

		// TODO: Idea is to transform the data I guess? To be able
		// to have sort of an encoder that applies to the API from this
		// service.
		data = append(data, '!', '!', '!')

		// 
		// Encoding step...
		// 
		// 
		return data
	})
	stream.SubscribeFilter(func (data []byte){
		// Decoding step.
		
		// stream.c.NATS()
		log.Println("<><><><><>", data)
	})

	stream.Subscribe(func(m *nats.Msg) {
		log.Println("-->", m.Subject, m)
	})

	// Creates a service object internally.
	// These have to be an interface to allow multiple type of usage.
	nc.Service("help").Subscribe(func(m *nats.Msg) {
		log.Println("->>", m.Subject, m.Reply, m)

		// Respond skips the filters so can't use the object?
		// m.Respond([]byte("+OK"))
	})

	// Sends all subscriptions and does a flush.
	// If there is any error during connect, like an auth error
	// on a subscription then
	err := nc.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Publish([]byte("Hello World"))
	log.Println(err)

	//
	// Can't create new services after being connected? Why not?
	// HTTP Servers can't create new handlers after starting,
	// can create new subscriptions via subs though, so not really
	// constrained either.
	//

	// Can reuse the middleware that has been installed so far.
	svc := nc.Service("help")

	// Allocate 10 goroutines to process (allow being able to tune on-the-fly, or recreate subscription?)
	// svc := nc.Service("help").WithInflightGroup(10)

	resp, err := svc.Request([]byte("help please"), 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response:", resp)

	//
	// nc.Service("foo").WithContext(ctx).Request([]byte("bar"))
	select {}
}
