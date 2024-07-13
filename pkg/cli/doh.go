package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/miekg/dns"
)

func QueryDoH(server, domain, recordType string, verbose bool) {
	rt, ok := dns.StringToType[recordType]
	if !ok {
		log.Fatalf("Unknown record type: %s\n", recordType)
	}
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), rt)
	dnsQuery, err := m.Pack()
	if err != nil {
		log.Fatalf("Failed to pack DNS query: %v\n", err)
	}

	url := fmt.Sprintf("https://%s/dns-query", server)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dnsQuery))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/dns-message")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform DoH query: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("DoH query failed: %s\n", body)
	}

	respMsg := new(dns.Msg)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read DoH response: %v\n", err)
	}
	err = respMsg.Unpack(respBody)
	if err != nil {
		log.Fatalf("Failed to unpack DoH response: %v\n", err)
	}

	if verbose {
		fmt.Println(respMsg.String())
	} else {
		if len(respMsg.Answer) == 0 {
			fmt.Println("No records found")
			return
		}

		for _, ans := range respMsg.Answer {
			fmt.Println(ans.String())
		}
	}
}
