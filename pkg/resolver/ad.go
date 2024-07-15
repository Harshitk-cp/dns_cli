package resolver

import (
	"net"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var (
	adBlockList   = sync.Map{}
	blockIP       = net.ParseIP("0.0.0.0")
	adBlockListMu sync.RWMutex
)

func InitAdBlockList(list map[string]struct{}) {
	adBlockListMu.Lock()
	defer adBlockListMu.Unlock()
	for domain := range list {
		adBlockList.Store(strings.ToLower(strings.TrimRight(domain, ".")), struct{}{})
	}
}

func IsAdDomain(domain string) bool {
	normalizedDomain := strings.ToLower(strings.TrimSpace(domain))
	normalizedDomain = strings.TrimRight(normalizedDomain, ".")

	adBlockListMu.RLock()
	defer adBlockListMu.RUnlock()

	_, found := adBlockList.Load(normalizedDomain)
	return found
}

func blockResponse(question dns.Question, r *dns.Msg) *dns.Msg {
	response := new(dns.Msg)
	response.SetRcode(r, dns.RcodeSuccess)
	response.RecursionAvailable = true
	response.Authoritative = true
	response.Question = []dns.Question{question}

	createBlockedRR := func(name string, rrType uint16, addr net.IP) dns.RR {
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: rrType,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: addr,
		}
	}

	switch question.Qtype {
	case dns.TypeA:
		response.Answer = append(response.Answer, createBlockedRR(question.Name, dns.TypeA, blockIP))
	case dns.TypeAAAA:
		response.Answer = append(response.Answer, createBlockedRR(question.Name, dns.TypeAAAA, blockIP))
	default:
		response.SetRcode(r, dns.RcodeNameError)
	}

	return response
}
