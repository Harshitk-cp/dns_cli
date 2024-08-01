package cache

import (
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type ResolverFunc func(domain string) (*dns.Msg, error)

type DNSCache struct {
	lru           *LRUCache
	ticker        *time.Ticker
	stopCh        chan struct{}
	usageCounts   map[string]int
	usageCountsMu sync.RWMutex
	popular       []string
	resolver      ResolverFunc
	prefetcher    *Prefetcher
}

func NewDNSCache() *DNSCache {
	c := &DNSCache{
		lru:         NewLRUCache(5),
		ticker:      time.NewTicker(time.Minute),
		stopCh:      make(chan struct{}),
		usageCounts: make(map[string]int),
	}
	c.prefetcher = NewPrefetcher(c)
	c.prefetcher.Start()
	go c.cleanupExpiredEntries()
	return c
}

func (c *DNSCache) SetResolver(resolver ResolverFunc) {
	c.resolver = resolver
}

func (c *DNSCache) Get(key string) (*dns.Msg, bool) {
	log.Printf("Cache hit for %s", key)
	return c.lru.Get(key)
}

func (c *DNSCache) Set(key string, msg *dns.Msg, ttl time.Duration) {
	log.Printf("Cache set for %s", key)
	c.lru.Set(key, msg, ttl)

	c.usageCountsMu.Lock()
	defer c.usageCountsMu.Unlock()
	c.usageCounts[key]++
}

func (c *DNSCache) cleanupExpiredEntries() {
	for {
		select {
		case <-c.ticker.C:
			log.Println("Cleaning up expired entries")
		case <-c.stopCh:
			return
		}
	}
}

func (c *DNSCache) Stop() {
	close(c.stopCh)
	c.ticker.Stop()
}
