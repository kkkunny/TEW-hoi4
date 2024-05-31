package localisation

import (
	"fmt"
	"path/filepath"
	"testing"

	stlmaps "github.com/kkkunny/stl/container/maps"

	"github.com/kkkunny/TEW-hoi4/config"
)

func TestParseLocalisation(t *testing.T) {
	locs, err := ParseChineseLocalisationDir(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld"))
	if err != nil {
		panic(err)
	}
	_, loc := stlmaps.Random(locs)
	fmt.Printf(loc.Encode())
}
