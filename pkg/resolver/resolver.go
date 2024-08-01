package resolver

import (
	"fmt"
	"log"
	"time"

	"github.com/Harshitk-cp/dns-cli/pkg/cache"
	"github.com/miekg/dns"
)

var (
	rootServers []string
	dnsCache    *cache.DNSCache
)

func init() {
	rootServers = []string{
		"198.41.0.4", "199.9.14.201", "192.33.4.12", "199.7.91.13",
		"192.203.230.10", "192.5.5.241", "192.112.36.4", "198.97.190.53",
	}

	dnsCache = cache.NewDNSCache()
	dnsCache.SetResolver(resolveDomain)
}

func resolveDomain(domain string) (*dns.Msg, error) {
	log.Printf("resolving domain: %s", domain)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	in, err := ResolveDNS(m.Question[0], m)
	if err != nil {
		return nil, err
	}
	return in, nil
}

func ResolveDNS(question dns.Question, r *dns.Msg) (*dns.Msg, error) {
	key := question.String()
	if cachedMsg, found := dnsCache.Get(key); found {
		return cachedMsg, nil
	}

	if IsAdDomain(question.Name) {
		return blockResponse(question, r), nil
	}

	servers := rootServers
	for i := 0; i < 10; i++ {
		response, err := queryServers(question, servers)
		if err != nil {
			return nil, err
		}

		if response.Rcode == dns.RcodeSuccess && len(response.Answer) > 0 {
			ttl := time.Duration(response.Answer[0].Header().Ttl) * time.Second
			dnsCache.Set(question.String(), response, ttl)
			return response, nil
		}

		servers = extractNameservers(response)
		if len(servers) == 0 {
			return response, nil
		}
	}
	return nil, fmt.Errorf("resolution incomplete after 10 iterations")
}

func extractNameservers(msg *dns.Msg) []string {
	var servers []string
	for _, rr := range msg.Ns {
		if ns, ok := rr.(*dns.NS); ok {
			servers = append(servers, ns.Ns)
		}
	}
	for _, rr := range msg.Extra {
		if a, ok := rr.(*dns.A); ok {
			servers = append(servers, a.A.String())
		}
	}
	return servers
}

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	question := r.Question[0]

	var response *dns.Msg
	var err error

	if question.Qtype == dns.TypePTR {
		response, err = resolveReverseDNS(question, r)
	} else {
		response, err = ResolveDNS(question, r)
	}

	if err != nil {
		log.Printf("Failed to resolve: %s\n", err.Error())
		dns.HandleFailed(w, r)
		return
	}

	response.Id = r.Id
	w.WriteMsg(response)
}
