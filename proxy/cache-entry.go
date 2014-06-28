package main

import (
	"sync"
	"time"

	"github.com/apokalyptik/photostream/client"
)

type cacheEntry struct {
	lock   sync.RWMutex
	birth  time.Time
	client *photostream.Client
	stream *photostream.WebStream
	assets *photostream.Assets
}

func (c *cacheEntry) getStream() (*photostream.WebStream, error) {
	c.lock.RLock()
	if c.stream != nil {
		c.lock.RUnlock()
		return c.stream, nil
	}
	c.lock.RUnlock()
	c.lock.Lock()
	defer c.lock.Unlock()
	s, e := c.client.Feed()
	if e != nil {
		return nil, e
	}
	c.stream = s
	return s, nil
}
