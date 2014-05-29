package photostream

import "time"

type AssetLocation struct {
	Hosts  []string `mapstructure:"hosts"`
	Scheme string   `mapstructure:"scheme"`
	Assets *Assets
}

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
