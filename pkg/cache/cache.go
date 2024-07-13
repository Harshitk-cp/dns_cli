package cache

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type DNSCache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

type cacheEntry struct {
	msg       *dns.Msg
	expiresAt time.Time
}

func NewDNSCache() *DNSCache {
	return &DNSCache{store: make(map[string]cacheEntry)}
}

func (c *DNSCache) Get(key string) (*dns.Msg, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.store[key]
	if !found || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.msg, true
}

func (c *DNSCache) Set(key string, msg *dns.Msg, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = cacheEntry{msg: msg, expiresAt: time.Now().Add(ttl)}
}
