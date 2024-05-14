package history

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/kkkunny/TEW-hoi4/config"
)

func TestParseState(t *testing.T) {
	// state, err := ParseStateDir(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld"))
	// if err != nil {
	// 	panic(err)
	// }
	// data, err := json.MarshalIndent(state[:10], "  ", "")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(data))
	state, err := ParseState(filepath.Join(config.HOI4MyModPath, "TheEmptyWorld", "history", "states", "4885-STATE_4885.txt"))
	if err != nil {
		panic(err)
	}
	fmt.Println(state.Encode())
}
