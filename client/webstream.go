package photostream

import (
	"errors"
	"fmt"
	"time"
)

type WebStreamItemDerivative struct {
	Checksum string `mapstructure:"checksum"`
	Size     int    `mapstructure:"fileSize"`
	Height   int    `mapstructure:"height"`
	State    string `mapstructure:"state"`
	Width    int    `mapstructure:"width"`
	item     *WebStreamItem
}

func (w *WebStreamItemDerivative) GetURLs() ([]string, error) {
	assets, err := w.item.stream.GetAssets()
	if err != nil {
		return nil, err
	}
	if v, ok := assets.Items[w.Checksum]; ok {
		rval := make([]string, len(v.location.Hosts))
		for k, h := range v.location.Hosts {
			rval[k] = fmt.Sprintf("%s://%s%s", v.location.Scheme, h, v.Path)
		}
		return rval, nil
	}
	return nil, errors.New("asset not found for item/derivative")
}

type WebStreamItem struct {
	BatchDateCreated string                              `mapstructure:"batchDateCreated"`
	BatchGuid        string                              `mapstructure:"batchGuid"`
	Caption          string                              `mapstructure:"caption"`
	FullName         string                              `mapstructure:"contributorFullName"`
	FirstName        string                              `mapstructure:"contributorFirstName"`
	LastName         string                              `mapstructure:"contributorLastName"`
	Created          string                              `mapstructure:"dateCreated"`
	Derivatives      map[string]*WebStreamItemDerivative `mapstructure:"derivatives"`
	Type             string                              `mapstructure:"mediaAssetType"`
	GUID             string                              `mapstructure:"photoGuid"`
	stream           *WebStream
}

func (w *WebStreamItem) init() {
	for k, _ := range w.Derivatives {
		w.Derivatives[k].item = w
	}
}

type WebStream struct {
	Items         interface{}     `mapstructure:"items"`
	ItemsReturned int             `mapstructure:"itemsReturned"`
	Media         []WebStreamItem `mapstructure:"photos"`
	Ctag          string          `mapstructure:"streamCtag"`
	Name          string          `mapstructure:"streamName"`
	FirstName     string          `mapstructure:"userFirstName"`
	LastName      string          `mapstructure:"userLastName"`
	client        *Client
	assets        *Assets
}

func (w *WebStream) init() {
	for k, _ := range w.Media {
		w.Media[k].stream = w
		w.Media[k].init()
	}
}

func (w *WebStream) GetAssets() (*Assets, error) {
	if w.assets != nil {
		var v *AssetItem
		for _, v = range w.assets.Items {
			if v.expires.After(time.Now().Add(5 * (time.Duration(1) - time.Minute))) {
				return w.assets, nil
			}
			break
		}
	}
	forDerivatives := make([]*WebStreamItemDerivative, 0)
	for _, m := range w.Media {
		for _, d := range m.Derivatives {
			forDerivatives = append(forDerivatives, d)
		}
	}
	assets, err := w.client.Assets(w.Media)
	if err != nil {
		return nil, err
	}
	w.assets = assets
	return assets, nil
}
