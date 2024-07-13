package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

func QueryDoH(server, domain string, opts QueryOptions) error {
	rt, ok := dns.StringToType[opts.RecordType]
	if !ok {
		return fmt.Errorf("unknown record type: %s", opts.RecordType)
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), rt)
	dnsQuery, err := m.Pack()
	if err != nil {
		return fmt.Errorf("failed to pack DNS query: %w", err)
	}

	url := fmt.Sprintf("https://%s/dns-query", server)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dnsQuery))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/dns-message")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform DoH query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DoH query failed with status %d: %s", resp.StatusCode, body)
	}

	respMsg := new(dns.Msg)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read DoH response: %w", err)
	}
	err = respMsg.Unpack(respBody)
	if err != nil {
		return fmt.Errorf("failed to unpack DoH response: %w", err)
	}

	printResponse(respMsg, domain, opts)
	return nil
}
