package _map

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	stlslices "github.com/kkkunny/stl/container/slices"
)

type StateType string

const (
	StateTypeLand StateType = "land"
	StateTypeSea  StateType = "sea"
)

type StateDef struct {
	ID          int64     `json:"id"`
	StateType   StateType `json:"state_type"`
	Landform    string    `json:"landform"`
	ContinentID int64     `json:"continent_id"`
}

func ParseStateDef(path string) (map[int64]*StateDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	matches := regexp.MustCompile(`(\d+);(\d+);(\d+);(\d+);(sea|land);(false|true);(\w+);(\d+)`).FindAllStringSubmatch(string(data), -1)
	stateDefs, err := stlslices.MapError(matches, func(_ int, match []string) (*StateDef, error) {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, err
		}
		switch match[5] {
		case "land", "sea":
		default:
			return nil, fmt.Errorf("unknown state type `%s`", match[5])
		}
		continentID, err := strconv.ParseInt(match[8], 10, 64)
		if err != nil {
			return nil, err
		}
		return &StateDef{
			ID:          id,
			StateType:   StateType(match[5]),
			Landform:    match[7],
			ContinentID: continentID,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return stlslices.ToMap[*StateDef, []*StateDef, int64, *StateDef, map[int64]*StateDef](stateDefs, func(def *StateDef) (int64, *StateDef) {
		return def.ID, def
	}), nil
}

func (def *StateDef) Encode() string {
	return fmt.Sprintf("%d;%d;%d;%d;%s;%s;%s;%d", def.ID, 0, 0, 0, def.StateType, "false", def.Landform, def.ContinentID)
}
