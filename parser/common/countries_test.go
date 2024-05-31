package common

import (
	"fmt"
	"path/filepath"
	"testing"

	stlmaps "github.com/kkkunny/stl/container/maps"

	"github.com/kkkunny/TEW-hoi4/config"
)

func TestParseCountries(t *testing.T) {
	// country, err := ParseCountryDef(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld", "common", "countries", "African.txt"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(country.Encode())
	// countryColors, err := ParseCountryColors(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld", "common", "countries", "cosmetic.txt"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(stlslices.First(countryColors).Encode())
	countryDefs, _, err := ParseCountriesDir(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld"))
	if err != nil {
		panic(err)
	}
	_, v := stlmaps.Random(countryDefs)
	fmt.Println(v.Encode())
}
