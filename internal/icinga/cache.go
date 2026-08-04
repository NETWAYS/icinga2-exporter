package icinga

import (
	"sync"
	"time"
)

type entry struct {
	data      []byte
	fetchedAt time.Time
}

type Cache struct {
	m   sync.Map
	ttl time.Duration
}

// NewCache returns a Cache that caches entries for the given ttl.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl}
}

func (c *Cache) Set(key string, b []byte) {
	if c.ttl <= 0 {
		return
	}

	// We're not creating a copy here, data never is mutated
	c.m.Store(key, &entry{
		data:      b,
		fetchedAt: time.Now(),
	})
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if c.ttl <= 0 {
		return nil, false
	}

	if v, ok := c.m.Load(key); ok {
		e, ok := v.(*entry)

		if !ok {
			return nil, false
		}

		if c.ttl > 0 && time.Since(e.fetchedAt) > c.ttl {
			c.m.Delete(key)
			return nil, false
		}
		// Again we're not creating a copy here, data never is mutated
		return e.data, true
	}

	return nil, false
}
