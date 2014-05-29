package photostream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type js map[string]interface{}

type Client struct {
	key   string
	shard string
	base  string
}

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
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected response: %s", resp.Status))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           webstream,
	})
	err = decoder.Decode(feed)
	if err != nil {
		return nil, err
	}
	webstream.client = c
	webstream.init()
	return webstream, nil
}

func New(key string) *Client {
	shard, _ := strconv.ParseInt(key[1:2], 36, 0)
	return &Client{
		key:   key,
		shard: fmt.Sprintf("%02d", shard),
		base:  fmt.Sprintf("https://p%02d-sharedstreams.icloud.com/%s/sharedstreams", shard, key),
	}
}
