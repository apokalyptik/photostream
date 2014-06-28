package main

import (
	"time"

	"github.com/apokalyptik/photostream/client"
)

type request struct {
	stream   string
	response chan *cacheEntry
}

var fetch = make(chan request)
var cacheDuration = time.Duration(5 * time.Minute)

func mindEngine() {
	var cache = make(map[string]*cacheEntry)
	var cacheCleanTicker = time.Tick(time.Duration(90 * time.Second))
	for {
		select {
		case r := <-fetch:
			if v, ok := cache[r.stream]; ok {
				r.response <- v
			} else {
				nv := &cacheEntry{
					birth:  time.Now(),
					client: photostream.New(r.stream),
				}
				cache[r.stream] = nv
				r.response <- nv
			}
		case <-cacheCleanTicker:
			for k, v := range cache {
				if time.Since(v.birth) > cacheDuration {
					delete(cache, k)
				}
			}
		}
	}
}
