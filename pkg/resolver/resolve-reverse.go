package resolver

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func resolveReverseDNS(question dns.Question, r *dns.Msg) (*dns.Msg, error) {
	ip := extractIPFromReverseDomain(question.Name)

	names, err := net.LookupAddr(ip)
	if err != nil {

		response := new(dns.Msg)
		response.SetRcode(r, dns.RcodeNameError)
		response.Authoritative = true
		return response, nil
	}

	log.Printf("Resolving reverse DNS query for IP: %s\n", ip)

	response := new(dns.Msg)
	response.SetReply(new(dns.Msg))
	response.Question = r.Question
	response.Authoritative = true

	for _, name := range names {
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    600,
			},
			Ptr: dns.Fqdn(name),
		}
		response.Answer = append(response.Answer, rr)
	}

	return response, nil
}

func extractIPFromReverseDomain(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) < 5 || parts[len(parts)-3] != "in-addr" || parts[len(parts)-2] != "arpa" {
		return ""
	}

	ipParts := make([]string, 4)
	for i := 0; i < 4; i++ {
		ipParts[3-i] = parts[i]
	}

	return strings.Join(ipParts, ".")
}
