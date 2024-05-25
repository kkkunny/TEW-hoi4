package common

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	stlslices "github.com/kkkunny/stl/container/slices"
)

type CountryTag struct {
	ID         string `json:"id"`
	DefinePath string `json:"path"`
}

func ParseCountryTag(path string) ([]*CountryTag, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contentBuffer strings.Builder
	reader := bufio.NewReader(file)
	for {
		lineData, err := reader.ReadBytes('\n')
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if i := bytes.IndexByte(lineData, '#'); i >= 0 {
			lineData = lineData[:i]
		}
		_, err = contentBuffer.Write(lineData)
		if err != nil {
			return nil, err
		}
	}

	content := strings.ReplaceAll(contentBuffer.String(), "\n", "")
	matches := regexp.MustCompile(`(\w{3})\s*=\s*"(.+?)"`).FindAllStringSubmatch(content, -1)
	countryTags := stlslices.Map(matches, func(_ int, match []string) *CountryTag {
		return &CountryTag{
			ID:         match[1],
			DefinePath: match[2],
		}
	})
	return countryTags, nil
}

func (state *CountryTag) Encode() string {
	return state.ID + " = \"" + state.DefinePath + "\""
}

func ParseCountryTagsDir(modPath string) ([]*CountryTag, error) {
	tagInfos, err := os.ReadDir(filepath.Join(modPath, "common", "country_tags"))
	if err != nil {
		return nil, err
	}
	return stlslices.FlatMapError(tagInfos, func(_ int, e os.DirEntry) ([]*CountryTag, error) {
		if e.IsDir() {
			return nil, nil
		}
		tags, err := ParseCountryTag(filepath.Join(modPath, "common", "country_tags", e.Name()))
		if err != nil {
			return nil, fmt.Errorf("country_tag %s parse error: %s", e.Name(), err.Error())
		}
		return tags, nil
	})
}

func ParseCountryTagsDirNotDynamic(modPath string) ([]*CountryTag, error) {
	tags, err := ParseCountryTagsDir(modPath)
	if err != nil {
		return nil, err
	}
	tags = stlslices.Filter(tags, func(_ int, tag *CountryTag) bool {
		return !regexp.MustCompile(`D\d{2}`).MatchString(tag.ID)
	})
	return tags, nil
}
