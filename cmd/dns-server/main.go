package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Harshitk-cp/dns-cli/pkg/doh"
	"github.com/Harshitk-cp/dns-cli/pkg/server"
	"github.com/Harshitk-cp/dns-cli/pkg/utils"
)

func main() {
	config := utils.LoadConfig("config/config.yaml")
	log.Println("Starting DNS and DoH servers...")

	go server.StartDNSServer(config.DNSServerAddr)

	go doh.StartDoHServer(config.DoHCertFile, config.DoHKeyFile, config.DOHServerAddr)

	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	server.StopDNSServer()
	doh.StopDoHServer()
	log.Println("Servers stopped")
}
