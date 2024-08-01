package cache

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"io"
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type LRUCache struct {
	mu      sync.RWMutex
	store   map[string]*list.Element
	ll      *list.List
	maxSize int
}

type cacheEntry struct {
	msg       []byte
	expiresAt time.Time
}

type entry struct {
	key   string
	value cacheEntry
}

func NewLRUCache(maxSize int) *LRUCache {
	return &LRUCache{
		store:   make(map[string]*list.Element),
		ll:      list.New(),
		maxSize: maxSize,
	}
}

func (c *LRUCache) Get(key string) (*dns.Msg, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ele, found := c.store[key]; found {
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*entry)
		if time.Now().After(entry.value.expiresAt) {
			c.ll.Remove(ele)
			delete(c.store, key)
			return nil, false
		}
		msg, err := decompress(entry.value.msg)
		if err != nil {
			return nil, false
		}
		return msg, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, msg *dns.Msg, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	compressedMsg, err := compress(msg)
	if err != nil {
		return
	}

	if ele, found := c.store[key]; found {
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*entry)
		entry.value = cacheEntry{msg: compressedMsg, expiresAt: time.Now().Add(ttl)}
	} else {
		if c.ll.Len() >= c.maxSize {
			c.removeOldest()
		}
		entry := &entry{
			key:   key,
			value: cacheEntry{msg: compressedMsg, expiresAt: time.Now().Add(ttl)},
		}
		ele := c.ll.PushFront(entry)
		c.store[key] = ele
	}
}

func (c *LRUCache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		entry := ele.Value.(*entry)
		delete(c.store, entry.key)
		log.Printf("Cache evicted for %s", entry.key)
	}
}

func compress(msg *dns.Msg) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	defer w.Close()
	data, err := msg.Pack()
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompress(data []byte) (*dns.Msg, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	msgData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	msg := new(dns.Msg)
	err = msg.Unpack(msgData)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
