package history

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kkkunny/stl/container/linkedhashmap"
	stlslices "github.com/kkkunny/stl/container/slices"
	"golang.org/x/exp/maps"
)

type State struct {
	ID            int64              `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Manpower      int64              `json:"manpower,omitempty"`
	Category      string             `json:"state_category,omitempty"`
	Impassable    bool               `json:"impassable,omitempty"`
	LocalSupplies float64            `json:"local_supplies"`
	Resources     map[string]float64 `json:"resources,omitempty"`
	Provinces     []int64            `json:"provinces,omitempty"`
	History       struct {
		Owner             string                     `json:"owner,omitempty"`
		Cores             []string                   `json:"-"`
		Claims            []string                   `json:"-"`
		CommonBuildings   map[string]int64           `json:"buildings,omitempty"`
		ProvinceBuildings map[int64]map[string]int64 `json:"-"`
		VictoryPoints     []int64                    `json:"victory_points,omitempty"`
	} `json:"history,omitempty"`
}

func ParseState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{','})
	data = bytes.ReplaceAll(data, []byte{'\r'}, nil)
	data = regexp.MustCompile(`\{.*}`).Find(data)
	data = bytes.ReplaceAll(data, []byte{'"'}, []byte{' '})

	content := regexp.MustCompile(`(\w+)\b`).ReplaceAllString(string(data), "\"${1}\"")
	content = regexp.MustCompile(`"(\d+)"`).ReplaceAllString(content, "${1}")
	content = strings.ReplaceAll(content, " ", ",")
	content = strings.ReplaceAll(content, "\t", ",")
	content = regexp.MustCompile(`,+`).ReplaceAllString(content, ",")
	content = strings.ReplaceAll(content, "{,", "{")
	content = strings.ReplaceAll(content, ",}", "}")
	content = strings.ReplaceAll(content, "\"yes\"", "true")
	content = strings.ReplaceAll(content, "\"no\"", "false")
	content = regexp.MustCompile(`,?=,?`).ReplaceAllString(content, ":")

	coreResults := regexp.MustCompile(`"add_core_of":"(.+?)"`).FindAllStringSubmatch(content, -1)
	cores := stlslices.Map(coreResults, func(_ int, e []string) string {
		return e[1]
	})
	content = regexp.MustCompile(`,?"add_core_of":"(.+?)"`).ReplaceAllString(content, "")

	claimResults := regexp.MustCompile(`"add_claim_of":"(.+?)"`).FindAllStringSubmatch(content, -1)
	claims := stlslices.Map(claimResults, func(_ int, e []string) string {
		return e[1]
	})
	content = regexp.MustCompile(`,?"add_claim_of":"(.+?)"`).ReplaceAllString(content, "")

	content = regexp.MustCompile(`\{([\d,]+)}`).ReplaceAllString(content, "[${1}]")
	content = regexp.MustCompile(`(\d+):`).ReplaceAllString(content, "\"${1}\":")

	provinceBuildingResults := regexp.MustCompile(`"(\d+)":(\{.+?})`).FindAllStringSubmatch(content, -1)
	provinceBuildingList, err := stlslices.MapError(provinceBuildingResults, func(_ int, e []string) (map[int64]map[string]int64, error) {
		provinceID, err := strconv.ParseInt(e[1], 10, 64)
		if err != nil {
			return nil, err
		}
		buildings := make(map[string]int64)
		err = json.Unmarshal([]byte(e[2]), &buildings)
		if err != nil {
			return nil, err
		}
		return map[int64]map[string]int64{provinceID: buildings}, nil
	})
	provinceBuildings := make(map[int64]map[string]int64, len(provinceBuildingList))
	for _, provinceBuilding := range provinceBuildingList {
		k := stlslices.First(maps.Keys(provinceBuilding))
		provinceBuildings[k] = provinceBuilding[k]
	}
	if err != nil {
		return nil, err
	}
	content = regexp.MustCompile(`,?"(\d+)":(\{.+?})`).ReplaceAllString(content, "")

	content = strings.ReplaceAll(content, "{,", "{")
	content = strings.ReplaceAll(content, ",}", "}")

	var state State
	err = json.Unmarshal([]byte(content), &state)
	if err != nil {
		return nil, err
	}
	state.History.Cores = cores
	state.History.Claims = claims
	state.History.ProvinceBuildings = provinceBuildings
	return &state, nil
}

type innerState struct {
	ID            int64                                    `json:"id,omitempty"`
	Name          string                                   `json:"name,omitempty"`
	Manpower      int64                                    `json:"manpower,omitempty"`
	Category      string                                   `json:"state_category,omitempty"`
	Impassable    bool                                     `json:"impassable,omitempty"`
	LocalSupplies float64                                  `json:"local_supplies"`
	Resources     map[string]float64                       `json:"resources,omitempty"`
	Provinces     []int64                                  `json:"provinces,omitempty"`
	History       linkedhashmap.LinkedHashMap[string, any] `json:"history,omitempty"`
}

func (state *State) Encode() string {
	cache := &innerState{
		ID:            state.ID,
		Name:          state.Name,
		Manpower:      state.Manpower,
		Category:      state.Category,
		Impassable:    state.Impassable,
		LocalSupplies: state.LocalSupplies,
		Resources:     state.Resources,
		Provinces:     state.Provinces,
		History:       linkedhashmap.NewLinkedHashMap[string, any](),
	}
	cache.History.Set("owner", state.History.Owner)
	if len(state.History.Cores) != 0 {
		for i, core := range state.History.Cores {
			cache.History.Set(fmt.Sprintf("add_core_of_%d", i), core)
		}
	}
	if len(state.History.Claims) != 0 {
		for i, claim := range state.History.Claims {
			cache.History.Set(fmt.Sprintf("add_claim_of_%d", i), claim)
		}
	}
	vicPoints := stlslices.Filter(stlslices.Map(state.History.VictoryPoints, func(i int, e int64) string {
		if i%2 != 0 {
			return ""
		}
		return fmt.Sprintf("%d %d", state.History.VictoryPoints[i], state.History.VictoryPoints[i+1])
	}), func(_ int, e string) bool {
		return e != ""
	})
	if len(vicPoints) != 0 {
		cache.History.Set("victory_points", vicPoints)
	}
	if len(state.History.CommonBuildings) != 0 {
		cache.History.Set("buildings", state.History.CommonBuildings)
	}
	if len(state.History.ProvinceBuildings) != 0 {
		for provinceID, building := range state.History.ProvinceBuildings {
			cache.History.Set(strconv.FormatInt(provinceID, 10), building)
		}
	}

	data, _ := json.MarshalIndent(cache, "", "  ")
	content := strings.ReplaceAll(string(data), ":", " =")
	content = strings.ReplaceAll(content, ",", "")
	content = strings.ReplaceAll(content, " true", " yes")
	content = strings.ReplaceAll(content, " false", " no")
	content = strings.ReplaceAll(content, "\"", "")
	content = regexp.MustCompile(`name = (.+)`).ReplaceAllString(content, "name = \"${1}\"")
	content = regexp.MustCompile(`add_core_of_\d+`).ReplaceAllString(content, "add_core_of")
	content = regexp.MustCompile(`add_claim_of_\d+`).ReplaceAllString(content, "add_claim_of")
	content = strings.ReplaceAll(content, "[", "{")
	content = strings.ReplaceAll(content, "]", "}")
	return "state = " + content
}

func ParseStateDir(modPath string) ([]*State, error) {
	stateInfos, err := os.ReadDir(filepath.Join(modPath, "history", "states"))
	if err != nil {
		return nil, err
	}
	return stlslices.FlatMapError(stateInfos, func(_ int, e os.DirEntry) ([]*State, error) {
		if e.IsDir() {
			return nil, nil
		}
		state, err := ParseState(filepath.Join(modPath, "history", "states", e.Name()))
		if err != nil {
			return nil, fmt.Errorf("state %s parse error: %s", e.Name(), err.Error())
		}
		return []*State{state}, nil
	})
}
