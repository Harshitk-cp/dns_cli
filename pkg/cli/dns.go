package cli

import (
	"fmt"

	"github.com/miekg/dns"
)

func QueryDNS(server, domain string, opts QueryOptions) error {
	c := new(dns.Client)
	m := new(dns.Msg)
	rt, ok := dns.StringToType[opts.RecordType]
	if !ok {
		return fmt.Errorf("unknown record type: %s", opts.RecordType)
	}
	m.SetQuestion(dns.Fqdn(domain), rt)
	r, _, err := c.Exchange(m, server+":53")
	if err != nil {
		return fmt.Errorf("failed to query DNS: %w", err)
	}

	printResponse(r, domain, opts)
	return nil
}

func printResponse(r *dns.Msg, domain string, opts QueryOptions) {
	if opts.Verbose {
		fmt.Println(r.String())
		return
	}

	if len(r.Answer) == 0 {
		fmt.Println("No records found")
		return
	}

	for _, ans := range r.Answer {
		if opts.RecordType == "PTR" {
			if ptr, ok := ans.(*dns.PTR); ok {
				fmt.Printf("%s => %s\n", domain, ptr.Ptr)
			}
		} else {
			fmt.Println(ans.String())
		}
	}
}
