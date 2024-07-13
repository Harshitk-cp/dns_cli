package cli

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func QueryDNS(server, domain, recordType string, verbose bool) {
	c := new(dns.Client)
	m := new(dns.Msg)
	rt, ok := dns.StringToType[recordType]
	if !ok {
		log.Fatalf("Unknown record type: %s\n", recordType)
	}
	m.SetQuestion(dns.Fqdn(domain), rt)
	r, _, err := c.Exchange(m, server+":53")
	if err != nil {
		log.Fatalf("Failed to query DNS: %v\n", err)
	}

	if verbose {
		fmt.Println(r.String())
	} else {
		if len(r.Answer) == 0 {
			fmt.Println("No records found")
			return
		}

		for _, ans := range r.Answer {
			fmt.Println(ans.String())
		}
	}
}
