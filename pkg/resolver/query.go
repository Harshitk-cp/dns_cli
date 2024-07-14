package resolver

import (
	"fmt"

	"github.com/miekg/dns"
)

func queryServers(question dns.Question, servers []string) (*dns.Msg, error) {
	for _, server := range servers {
		response, err := queryServer(question, server)
		if err == nil {
			for _, answer := range response.Answer {
				if cname, ok := answer.(*dns.CNAME); ok {
					cnameQuestion := dns.Question{
						Name:   cname.Target,
						Qtype:  question.Qtype,
						Qclass: dns.ClassINET,
					}
					cnameResponse, err := ResolveDNS(cnameQuestion)
					if err != nil {
						return nil, err
					}
					response.Answer = append(response.Answer, cnameResponse.Answer...)
				}
			}
			return response, nil
		}
	}
	return nil, fmt.Errorf("failed to query any server")
}

func queryServer(question dns.Question, server string) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(question.Name, question.Qtype)
	m.RecursionDesired = true

	client := new(dns.Client)
	in, _, err := client.Exchange(m, server+":53")
	if err != nil {
		return nil, err
	}

	if in.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("failed to get valid answer")
	}

	return in, nil
}
