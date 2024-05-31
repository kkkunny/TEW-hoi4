package localisation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kkkunny/stl/container/optional"
	stlslices "github.com/kkkunny/stl/container/slices"
)

type Localisation struct {
	Key   string                 `json:"key"`
	Index optional.Optional[int] `json:"index,omitempty"`
	Value string                 `json:"value"`
}

func ParseLocalisation(path string) (map[string]*Localisation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	content := string(data)
	matches := regexp.MustCompile(`(.+?)\s*:\s*(\d*)\s*"(.*)"`).FindAllStringSubmatch(content, -1)
	locs, err := stlslices.MapError(matches, func(_ int, match []string) (*Localisation, error) {
		var index optional.Optional[int]
		if match[2] != "" {
			indexVal, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				return nil, err
			}
			index = optional.Some(int(indexVal))
		}
		return &Localisation{
			Key:   match[1],
			Index: index,
			Value: match[3],
		}, nil
	})
	return stlslices.ToMap[*Localisation, []*Localisation, string, *Localisation, map[string]*Localisation](locs, func(loc *Localisation) (string, *Localisation) {
		return loc.Key, loc
	}), nil
}

func (loc *Localisation) Encode() string {
	if loc.Index.IsNone() {
		return fmt.Sprintf("%s: \"%s\"", loc.Key, loc.Value)
	} else {
		return fmt.Sprintf("%s:%d \"%s\"", loc.Key, loc.Index.MustValue(), loc.Value)
	}
}

func ParseLocalisationDir(modPath string) (map[string]map[string]*Localisation, error) {
	dirPath := filepath.Join(modPath, "localisation")
	locDirInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	locs := make(map[string]map[string]*Localisation)
	for _, locDirInfo := range locDirInfos {
		if !locDirInfo.IsDir() {
			continue
		}
		locs[locDirInfo.Name()] = make(map[string]*Localisation)

		locDirPath := filepath.Join(dirPath, locDirInfo.Name())
		locInfos, err := os.ReadDir(locDirPath)
		if err != nil {
			return nil, err
		}
		for _, locInfo := range locInfos {
			if locInfo.IsDir() || !strings.HasSuffix(locInfo.Name(), ".yml") {
				continue
			}
			fp := filepath.Join(locDirPath, locInfo.Name())
			oneLocs, err := ParseLocalisation(fp)
			if err != nil {
				return nil, fmt.Errorf("`%s` parse error: %s", fp, err.Error())
			}
			for k, oneLoc := range oneLocs {
				locs[locDirInfo.Name()][k] = oneLoc
			}
		}
	}
	return locs, nil
}

func ParseChineseLocalisationDir(modPath string) (map[string]*Localisation, error) {
	dirPath := filepath.Join(modPath, "localisation", "simp_chinese")
	locInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	locs := make(map[string]*Localisation)
	for _, locInfo := range locInfos {
		if locInfo.IsDir() || !strings.HasSuffix(locInfo.Name(), ".yml") {
			continue
		}
		fp := filepath.Join(dirPath, locInfo.Name())
		oneLocs, err := ParseLocalisation(fp)
		if err != nil {
			return nil, fmt.Errorf("`%s` parse error: %s", fp, err.Error())
		}
		for k, oneLoc := range oneLocs {
			locs[k] = oneLoc
		}
	}
	return locs, nil
}
