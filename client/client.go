// Package photostream impliments a client capable of reading
// all the data from a photostream publicly shared from an
// iCloud account
package photostream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type js map[string]interface{}

// Client is the base from which all your use of this api will
// normally stem. You should not create your own Client, and
// instead should use the New() function
type Client struct {
	key   string
	shard string
	base  string
}

// Assets if used to retrieve the asset information for a
// slice of WebStreamItem structs, or error if something went
// wrong.  It's likely very rare that you would use this
// directly unless you're doing something more advanced than
// normal.
func (c *Client) Assets(items []WebStreamItem) (*Assets, error) {
	assets := &Assets{}
	checksums := make([]string, len(items))
	for k, v := range items {
		checksums[k] = v.GUID
	}
	data, err := json.Marshal(map[string][]string{
		"photoGuids": checksums,
	})
	if err != nil {
		return nil, err
	}
	feed := js{}
	resp, err := http.Post(
		fmt.Sprintf("%s/webasseturls", c.base),
		"text/plain",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &feed)
	if err != nil {
		return nil, err
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           assets,
	})
	err = decoder.Decode(feed)
	if err != nil {
		return nil, err
	}
	assets.init()
	return assets, nil
}

// Feed pulls a list of items from the photostream, or returns
// an error. This is probably the first thing you use after
// getting youself a new Client.
func (c *Client) Feed() (*WebStream, error) {
	feed := js{}
	webstream := &WebStream{}
	resp, err := http.Post(
		fmt.Sprintf("%s/webstream", c.base),
		"text/plain",
		bytes.NewBufferString("{\"streamCtag\":null}"),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(1)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status)
	}
	fmt.Println(2)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(3)

	err = json.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	fmt.Println(4)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           webstream,
	})
	fmt.Println(5)
	err = decoder.Decode(feed)
	if err != nil {
		return nil, err
	}
	fmt.Println(6)
	webstream.client = c
	webstream.init()
	return webstream, nil
}

// New returns a new Client struct with the necessary information
// embedded into it to preform the agreed upon duties.
func New(key string) *Client {
	shard, _ := strconv.ParseInt(key[1:2], 36, 0)
	return &Client{
		key:   key,
		shard: fmt.Sprintf("%02d", shard),
		base:  fmt.Sprintf("https://p%02d-sharedstreams.icloud.com/%s/sharedstreams", shard, key),
	}
}
