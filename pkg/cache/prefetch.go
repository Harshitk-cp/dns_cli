package cache

import (
	"log"
	"sort"
	"strings"
	"time"
)

type Prefetcher struct {
	cache *DNSCache
}

func NewPrefetcher(cache *DNSCache) *Prefetcher {
	return &Prefetcher{
		cache: cache,
	}
}

func (p *Prefetcher) Start() {
	go p.monitorPopularDomains()
}

func (p *Prefetcher) monitorPopularDomains() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.updatePopularDomains()
		case <-p.cache.stopCh:
			return
		}
	}
}

func (p *Prefetcher) updatePopularDomains() {
	p.cache.usageCountsMu.RLock()
	defer p.cache.usageCountsMu.RUnlock()

	type domainCount struct {
		domain string
		count  int
	}
	var domains []domainCount
	for domain, count := range p.cache.usageCounts {
		domains = append(domains, domainCount{domain, count})
	}

	sort.Slice(domains, func(i, j int) bool {
		return domains[i].count > domains[j].count
	})

	topN := 5
	p.cache.popular = make([]string, 0, min(len(domains), topN))
	for i := 0; i < min(len(domains), topN); i++ {
		p.cache.popular = append(p.cache.popular, domains[i].domain)
	}

	p.prefetchPopular()
}

func (p *Prefetcher) prefetchPopular() {
	log.Printf("Prefetching popular domains: %v", p.cache.popular)
	for _, fullQuery := range p.cache.popular {
		domain := extractDomain(fullQuery)
		if _, found := p.cache.lru.Get(fullQuery); !found {
			msg, err := p.cache.resolver(domain)
			if err == nil && len(msg.Answer) > 0 {
				ttl := time.Duration(msg.Answer[0].Header().Ttl) * time.Second
				p.cache.Set(fullQuery, msg, ttl)
			} else {
				log.Printf("Failed to prefetch %s: %v", domain, err)
			}
		}
	}
}

func extractDomain(query string) string {
	query = strings.TrimLeft(query, "; ")
	parts := strings.Fields(query)
	if len(parts) > 0 {
		return strings.TrimSuffix(parts[0], ".")
	}
	return query
}
