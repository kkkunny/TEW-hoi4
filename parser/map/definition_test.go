package _map

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/kkkunny/TEW-hoi4/config"
)

func TestParseDefinition(t *testing.T) {
	defs, err := ParseStateDef(filepath.Join(config.TEWRootPath, "map", "definition.csv"))
	if err != nil {
		panic(err)
	}
	for _, def := range defs {
		fmt.Println(def.Encode())
	}
}
