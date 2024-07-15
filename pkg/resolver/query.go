package resolver

import (
	"fmt"

	"github.com/miekg/dns"
)

func queryServers(question dns.Question, servers []string) (*dns.Msg, error) {
	seen := make(map[string]struct{})
	for _, server := range servers {
		response, err := queryServer(question, server)
		if err == nil {
			var uniqueAnswers []dns.RR
			for _, answer := range response.Answer {
				if _, found := seen[answer.String()]; !found {
					seen[answer.String()] = struct{}{}
					uniqueAnswers = append(uniqueAnswers, answer)
				}
				if cname, ok := answer.(*dns.CNAME); ok {
					cnameQuestion := dns.Question{
						Name:   cname.Target,
						Qtype:  question.Qtype,
						Qclass: dns.ClassINET,
					}
					cnameResponse, err := ResolveDNS(cnameQuestion, response)
					if err != nil {
						return nil, err
					}
					for _, cnameAnswer := range cnameResponse.Answer {
						if _, found := seen[cnameAnswer.String()]; !found {
							seen[cnameAnswer.String()] = struct{}{}
							uniqueAnswers = append(uniqueAnswers, cnameAnswer)
						}
					}
				}
			}
			response.Answer = uniqueAnswers
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
