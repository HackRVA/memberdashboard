package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
)

type Resource struct {
	Address string
	Label   string
}

// listenForBroadcasts - devices will broadcast that they exist.
//  when a device is broadcasting, we tell the server that it is an available
//  resource to provision
func listenForBroadcasts() {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s sent this: %s\n", addr, buf[:n])
}

func registerResource() {

}

func forwardToServer() {

}

func forwardToResource() {

}
func browseForServices() {
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println(entry)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}

func setupMDNS() {
	// // Setup our service export
	// // host, _ := os.Hostname()
	// info := []string{"My awesome service"}
	// service, _ := mdns.NewMDNSService("resourceBRIDGE", "_foobar._tcp", "", "", 8000, nil, info)

	// // Create the mDNS server, defer shutdown
	// server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	// defer server.Shutdown()

	// // Make a channel for results and start listening
	// entriesCh := make(chan *mdns.ServiceEntry, 4)
	// entriesCh <- &mdns.ServiceEntry{}
	// go func() {
	// 	for entry := range entriesCh {
	// 		fmt.Printf("Got new entry: %v\n", entry)
	// 	}
	// }()

	// // Start the lookup
	// mdns.Lookup("_foobar._tcp", entriesCh)
	// close(entriesCh)

	server, err := zeroconf.Register("GoZeroconf", "_workstation._tcp", "local.", 42424, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		// Exit by timeout
	}

	log.Println("Shutting down.")
}

func lookupResource() {
}

// pollServer - check the server to see if we have any tasks that need doing
func pollServer() {

}

func main() {
	setupMDNS()
	lookupResource()
}
