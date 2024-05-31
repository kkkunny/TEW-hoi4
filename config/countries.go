package config

import (
	_ "embed"
	"encoding/json"

	"github.com/kkkunny/stl/container/optional"
	stlslices "github.com/kkkunny/stl/container/slices"
)

type Country struct {
	ID     string                      `json:"id"`
	Name   string                      `json:"name"`
	Region string                      `json:"region"`
	Color  optional.Optional[[3]uint8] `json:"color,omitempty"`
}

//go:embed countries.json
var countriesData []byte
var Countries = func() map[string]*Country {
	var countries []*Country
	err := json.Unmarshal(countriesData, &countries)
	if err != nil {
		panic(err)
	}
	return stlslices.ToMap[*Country, []*Country, string, *Country, map[string]*Country](countries, func(c *Country) (string, *Country) {
		return c.ID, c
	})
}()
