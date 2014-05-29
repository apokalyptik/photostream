package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/apokalyptik/photostream/client"
)

var streamKey string

func init() {
	flag.StringVar(&streamKey, "stream", "", "the photostream key, eg: A2GI9HKKGiWkZH")
	flag.Parse()
}

func main() {
	c := photostream.New(streamKey)
	feed, err := c.Feed()
	if err != nil {
		log.Fatalf("error fetching feed: %s", err.Error())
	}
	for _, item := range feed.Media {
		for _, v := range item.Derivatives {
			u, err := v.GetURLs()
			if err != nil {
				log.Fatalf("error finding derivative URL: %s", err.Error())
			}
			for _, url := range u {
				fmt.Println(url)
			}
		}
	}
}
