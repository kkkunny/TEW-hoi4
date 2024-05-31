package common

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	stlbasic "github.com/kkkunny/stl/basic"
	"github.com/kkkunny/stl/container/optional"
	stlslices "github.com/kkkunny/stl/container/slices"

	"github.com/kkkunny/TEW-hoi4/util"
)

type CountryDef struct {
	GraphicalCulture   optional.Optional[string]      `json:"graphical_culture,omitempty"`
	GraphicalCulture2D optional.Optional[string]      `json:"graphical_culture_2d,omitempty"`
	Color              optional.Optional[color.Color] `json:"color,omitempty"`
}

func ParseCountryDef(path string) (*CountryDef, error) {
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
	graphicalCulture := stlslices.Last(regexp.MustCompile(`graphical_culture\s*=\s*(.+)`).FindStringSubmatch(content))
	graphicalCulture2D := stlslices.Last(regexp.MustCompile(`graphical_culture_2d\s*=\s*(.+)`).FindStringSubmatch(content))
	colorMatches := regexp.MustCompile(`color\s*=\s*(\w*?)\s*{\s*([\d.]+?)\s*([\d.]+?)\s*([\d.]+?)\s*}`).FindStringSubmatch(content)
	var colorVal optional.Optional[color.Color]
	if len(colorMatches) > 0 {
		v1, err := strconv.ParseFloat(colorMatches[2], 64)
		if err != nil {
			return nil, err
		}
		v2, err := strconv.ParseFloat(colorMatches[3], 64)
		if err != nil {
			return nil, err
		}
		v3, err := strconv.ParseFloat(colorMatches[4], 64)
		if err != nil {
			return nil, err
		}
		clr, err := util.NewColorByMode(colorMatches[1], v1, v2, v3)
		if err != nil {
			return nil, err
		}
		colorVal = optional.Some(clr)
	}
	return &CountryDef{
		GraphicalCulture:   stlbasic.Ternary(graphicalCulture == "", optional.None[string](), optional.Some(graphicalCulture)),
		GraphicalCulture2D: stlbasic.Ternary(graphicalCulture2D == "", optional.None[string](), optional.Some(graphicalCulture2D)),
		Color:              colorVal,
	}, nil
}

func (c *CountryDef) Encode() string {
	var buf strings.Builder

	if c.GraphicalCulture.IsSome() {
		buf.WriteString(fmt.Sprintf("graphical_culture = %s\n", c.GraphicalCulture.MustValue()))
	}
	if c.GraphicalCulture2D.IsSome() {
		buf.WriteString(fmt.Sprintf("graphical_culture_2d = %s\n\n", c.GraphicalCulture2D.MustValue()))
	}
	if c.Color.IsSome() {
		r, g, b := util.GetRGB(c.Color.MustValue())
		buf.WriteString(fmt.Sprintf("color = rgb { %d %d %d }", r, g, b))
	}

	return buf.String()
}

type CountryColor struct {
	Country string      `json:"country"`
	Color   color.Color `json:"color"`
	ColorUI color.Color `json:"color_ui"`
}

func ParseCountryColors(path string) ([]*CountryColor, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	content := strings.ReplaceAll(string(data), "\n", " ")
	countryColorMatchesList := regexp.MustCompile(`(\w+)\s*=\s*{\s*color\s*=\s*(\w*?)\s*{\s*([\d.]+?)\s*([\d.]+?)\s*([\d.]+?)\s*}\s*color_ui\s*=\s*(\w*?)\s*{\s*([\d.]+?)\s*([\d.]+?)\s*([\d.]+?)\s*}`).FindAllStringSubmatch(content, -1)
	return stlslices.MapError(countryColorMatchesList, func(_ int, countryColorMatches []string) (*CountryColor, error) {
		v1, err := strconv.ParseFloat(countryColorMatches[3], 64)
		if err != nil {
			return nil, err
		}
		v2, err := strconv.ParseFloat(countryColorMatches[4], 64)
		if err != nil {
			return nil, err
		}
		v3, err := strconv.ParseFloat(countryColorMatches[5], 64)
		if err != nil {
			return nil, err
		}
		clr, err := util.NewColorByMode(countryColorMatches[2], v1, v2, v3)
		if err != nil {
			return nil, err
		}
		v1, err = strconv.ParseFloat(countryColorMatches[7], 64)
		if err != nil {
			return nil, err
		}
		v2, err = strconv.ParseFloat(countryColorMatches[8], 64)
		if err != nil {
			return nil, err
		}
		v3, err = strconv.ParseFloat(countryColorMatches[9], 64)
		if err != nil {
			return nil, err
		}
		clrUI, err := util.NewColorByMode(countryColorMatches[6], v1, v2, v3)
		if err != nil {
			return nil, err
		}
		return &CountryColor{
			Country: countryColorMatches[1],
			Color:   clr,
			ColorUI: clrUI,
		}, nil
	})
}

func (cc *CountryColor) Encode() string {
	r1, g1, b1 := util.GetRGB(cc.Color)
	r2, g2, b2 := util.GetRGB(cc.ColorUI)
	return fmt.Sprintf("%s = {\n\tcolor = rgb { %d %d %d }\n\tcolor_ui = rgb { %d %d %d }\n}", cc.Country, r1, g1, b1, r2, g2, b2)
}

func ParseCountriesDir(modPath string) (map[string]*CountryDef, map[string]*CountryColor, error) {
	countryInfos, err := os.ReadDir(filepath.Join(modPath, "common", "countries"))
	if err != nil {
		return nil, nil, err
	}
	countryDefs := make(map[string]*CountryDef)
	colors := make(map[string]*CountryColor)
	for _, countryInfo := range countryInfos {
		if countryInfo.IsDir() {
			continue
		}
		fp := filepath.Join(modPath, "common", "countries", countryInfo.Name())
		switch {
		case stlslices.Contain([]string{"colors.txt", "cosmetic.txt"}, countryInfo.Name()):
			ccs, err := ParseCountryColors(fp)
			if err != nil {
				return nil, nil, fmt.Errorf("`%s` parse error: %s", fp, err.Error())
			}
			for _, cc := range ccs {
				colors[cc.Country] = cc
			}
		case strings.HasSuffix(countryInfo.Name(), ".txt"):
			countryDef, err := ParseCountryDef(fp)
			if err != nil {
				return nil, nil, fmt.Errorf("`%s` parse error: %s", fp, err.Error())
			}
			countryDefs[strings.TrimSuffix(countryInfo.Name(), ".txt")] = countryDef
		}
	}
	return countryDefs, colors, nil
}
