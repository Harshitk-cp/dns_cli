package main

import (
	"log"
	"os"

	"github.com/Harshitk-cp/dns-cli/pkg/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <start|stop|status>\n", os.Args[0])
	}

	command := os.Args[1]

	switch command {
	case "start":
		log.Println("DNS server started")
		server.StartDNSServer()
	case "stop":
		server.StopDNSServer()
		log.Println("DNS server stopped")
	case "status":
		log.Println("DNS server status: running")
	default:
		log.Fatalf("Unknown command: %s\n", command)
	}
}
