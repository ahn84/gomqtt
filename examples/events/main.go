package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logrusorgru/aurora"

	mqtt "github.com/ahn84/gomqtt/server"
	"github.com/ahn84/gomqtt/server/events"
	"github.com/ahn84/gomqtt/server/listeners"
	"github.com/ahn84/gomqtt/server/listeners/auth"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("TCP"))

	server := mqtt.NewServer(nil)
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Add OnConnect Event Hook
	server.Events.OnConnect = func(cl events.ClientLike, pk events.Packet) {
		fmt.Printf("<< OnConnect client connected %s: %+v\n", cl.ID, pk)
	}

	// Add OnDisconnect Event Hook
	server.Events.OnDisconnect = func(cl events.ClientLike, err error) {
		fmt.Printf("<< OnDisconnect client disconnected %s: %v\n", cl.ID, err)
	}

	// Add OnSubscribe Event Hook
	server.Events.OnSubscribe = func(filter string, cl events.ClientLike, qos byte) {
		fmt.Printf("<< OnSubscribe client subscribed %s: %s %v\n", cl.ID, filter, qos)
	}

	// Add OnUnsubscribe Event Hook
	server.Events.OnUnsubscribe = func(filter string, cl events.ClientLike) {
		fmt.Printf("<< OnUnsubscribe client unsubscribed %s: %s\n", cl.ID, filter)
	}

	// Add OnMessage Event Hook
	server.Events.OnMessage = func(cl events.ClientLike, pk events.Packet) (pkx events.Packet, err error) {
		pkx = pk
		if string(pk.Payload) == "hello" {
			pkx.Payload = []byte("hello world")
			fmt.Printf("< OnMessage modified message from client %s: %s\n", cl.ID, string(pkx.Payload))
		} else {
			fmt.Printf("< OnMessage received message from client %s: %s\n", cl.ID, string(pkx.Payload))
		}

		// Example of using AllowClients to selectively deliver/drop messages.
		// Only a client with the id of `allowed-client` will received messages on the topic.
		if pkx.TopicName == "a/b/restricted" {
			pkx.AllowClients = []string{"allowed-client"} // slice of known client ids
		}

		return pkx, nil
	}

	// Demonstration of directly publishing messages to a topic via the
	// `server.Publish` method. Subscribe to `direct/publish` using your
	// MQTT client to see the messages.
	go func() {
		for range time.Tick(time.Second * 10) {
			server.Publish("direct/publish", []byte("scheduled message"), false)
			fmt.Println("> issued direct message to direct/publish")
		}
	}()

	fmt.Println(aurora.BgMagenta("  Started!  "))

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))
}
