package photostream

import (
	"errors"
	"fmt"
	"time"
)

// WebStreamItemDerivative represents the one of the sizes
// (derivatives) available for this particular media item
type WebStreamItemDerivative struct {
	Checksum string `mapstructure:"checksum"`
	Size     int    `mapstructure:"fileSize"`
	Height   int    `mapstructure:"height"`
	State    string `mapstructure:"state"`
	Width    int    `mapstructure:"width"`
	item     *WebStreamItem
}

// GetURLs returns all of the possible urls against which a GET
// request can be made to fetch the particular derivitive
// representation (that is, it gives you URLs to fetch this
// particular size with)
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

// WebStreamItem represents the data for a single media item
// in the photostream
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
	for k := range w.Derivatives {
		w.Derivatives[k].item = w
	}
}

// WebStream represents the total information returned in an
// initial request to a photostream.
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
	for k := range w.Media {
		w.Media[k].stream = w
		w.Media[k].init()
	}
}

// GetAssets makes a request to get all of the asset formation
// associated with this webstream data.  These assets are cacheable
// for a certain period of time, and so this function will give you
// a cached version upon calling the function again.  When the data
// is about to expire (in the next 5 minutes) it will make a new
// request for fresh information.  WebStreamItemDerivative.GetURLs()
// uses this internally
func (w *WebStream) GetAssets() (*Assets, error) {
	var forDerivatives []*WebStreamItemDerivative
	if w.assets != nil {
		var v *AssetItem
		for _, v = range w.assets.Items {
			if v.expires.After(time.Now().Add(5 * (time.Duration(1) - time.Minute))) {
				return w.assets, nil
			}
			break
		}
	}
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
