package resolver

import (
	"fmt"

	"github.com/miekg/dns"
)

func queryServers(question dns.Question, servers []string) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(question.Name, question.Qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	for _, server := range servers {
		response, _, err := c.Exchange(m, server+":53")
		if err == nil {
			return response, nil
		}
	}
	return nil, fmt.Errorf("failed to query all servers")
}
