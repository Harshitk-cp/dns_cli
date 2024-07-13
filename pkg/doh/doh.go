package doh

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Harshitk-cp/dns-cli/pkg/resolver"
	"github.com/miekg/dns"
)

func StartDoHServer(certFile, keyFile, addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/dns-query", handleDoHRequest)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Printf("Starting DoH server on %s\n", addr)
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to start DoH server: %s\n", err.Error())
	}
}

func handleDoHRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	var msg []byte
	var err error

	if r.Method == http.MethodGet {
		dnsQuery := r.URL.Query().Get("dns")
		if dnsQuery == "" {
			http.Error(w, "Missing DNS query", http.StatusBadRequest)
			return
		}
		msg, err = base64.RawURLEncoding.DecodeString(dnsQuery)
	} else if r.Method == http.MethodPost {
		msg, err = io.ReadAll(r.Body)
		defer r.Body.Close()
	}

	if err != nil {
		http.Error(w, "Invalid DNS query", http.StatusBadRequest)
		return
	}

	dnsMsg := new(dns.Msg)
	err = dnsMsg.Unpack(msg)
	if err != nil {
		http.Error(w, "Failed to parse DNS query", http.StatusBadRequest)
		return
	}

	response, err := resolver.ResolveDNS(dnsMsg.Question[0])
	if err != nil {
		http.Error(w, "Failed to resolve DNS query", http.StatusInternalServerError)
		return
	}

	response.Id = dnsMsg.Id
	responseBytes, err := response.Pack()
	if err != nil {
		http.Error(w, "Failed to pack DNS response", http.StatusInternalServerError)
		return
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/dns-json" {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/dns-json")
		w.Write(jsonResponse)
	} else {
		w.Header().Set("Content-Type", "application/dns-message")
		w.Write(responseBytes)
	}
}
