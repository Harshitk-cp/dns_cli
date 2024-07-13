package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Harshitk-cp/dns-cli/pkg/resolver"
	"github.com/miekg/dns"
)

var (
	dnsServer *dns.Server
)

func StartDNSServer() {
	dns.HandleFunc(".", resolver.HandleDNSRequest)
	dnsServer = &dns.Server{Addr: ":53", Net: "udp"}
	go func() {
		log.Println("Starting DNS server on :53")
		err := dnsServer.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start DNS server: %s\n", err.Error())
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	StopDNSServer()
}

func StopDNSServer() {
	if dnsServer != nil {
		dnsServer.Shutdown()
		log.Println("DNS server stopped")
	} else {
		log.Println("DNS server is not running")
	}
}
