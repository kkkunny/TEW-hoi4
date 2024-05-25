package common

import (
	"fmt"
	"path/filepath"
	"testing"

	stlbasic "github.com/kkkunny/stl/basic"

	"github.com/kkkunny/TEW-hoi4/config"
)

func TestParseCountryTag(t *testing.T) {
	// state, err := ParseStateDir(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld"))
	// if err != nil {
	// 	panic(err)
	// }
	// data, err := json.MarshalIndent(state[:10], "  ", "")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(data))
	tags, err := ParseCountryTagsDirNotDynamic(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld"))
	if err != nil {
		panic(err)
	}
	for i, tag := range tags {
		fmt.Printf("\t\t\t\t%s = {\n", stlbasic.Ternary(i == 0, "if", "else_if"))
		fmt.Printf("\t\t\t\t\tlimit = { original_tag = %s }\n", tag.ID)
		fmt.Printf("\t\t\t\t\tset_cosmetic_tag = %s_alias_%s\n", tag.ID, "anarchist")
		fmt.Println("\t\t\t\t}")
	}
}
