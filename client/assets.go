package photostream

import "time"

// AssetLocation describes a particular place against which
// asset GET requests can be made
type AssetLocation struct {
	Hosts  []string `mapstructure:"hosts"`
	Scheme string   `mapstructure:"scheme"`
	Assets *Assets
}

// AssetItem represents the information about a single media
// item
type AssetItem struct {
	Expires  string `mapstructure:"url_expiry"`
	Location string `mapstructure:"url_location"`
	Path     string `mapstructure:"url_path"`
	location *AssetLocation
	expires  time.Time
	assets   *Assets
}

func (i *AssetItem) init() {
	if v, ok := i.assets.Locations[i.Location]; ok {
		i.location = v
	}
	i.expires, _ = time.Parse(time.RFC3339, i.Expires)
}

// Assets represents the total data returned by a request for
// asset information
type Assets struct {
	Items     map[string]*AssetItem     `mapstructure:"items"`
	Locations map[string]*AssetLocation `mapstructure:"locations"`
}

func (a *Assets) init() {
	for _, v := range a.Items {
		v.assets = a
		v.init()
	}
}
