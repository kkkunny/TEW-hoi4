package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"math/rand/v2"
	"os"
	"path/filepath"

	stlbasic "github.com/kkkunny/stl/basic"
	"github.com/kkkunny/stl/container/hashset"
	"github.com/kkkunny/stl/container/optional"
	stlslices "github.com/kkkunny/stl/container/slices"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/exp/maps"

	"github.com/kkkunny/TEW-hoi4/config"
	"github.com/kkkunny/TEW-hoi4/parser/common"
	"github.com/kkkunny/TEW-hoi4/parser/localisation"
	"github.com/kkkunny/TEW-hoi4/util"
)

func RefreshCountries() error {
	modPath := config.TEWRootPath

	fmt.Println("生成国家tag文件中...")
	var countryTagBuffer bytes.Buffer
	for _, c := range config.Countries {
		ct := common.CountryTag{
			ID:         c.ID,
			DefinePath: "countries/" + c.Region + ".txt",
		}
		countryTagBuffer.WriteString(ct.Encode())
		countryTagBuffer.WriteString("\n")
	}
	err := os.WriteFile(filepath.Join(modPath, "common", "country_tags", "tew_auto_generate.txt"), countryTagBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成国家tag文件成功！")

	fmt.Println("生成国家名字文件中...")
	var countryLocBuffer bytes.Buffer
	countryLocBuffer.WriteString("l_simp_chinese:\n")
	for _, c := range config.Countries {
		loc := &localisation.Localisation{
			Key:   c.ID,
			Index: optional.Some(0),
			Value: c.Name,
		}
		countryLocBuffer.WriteByte(' ')
		countryLocBuffer.WriteString(loc.Encode())
		countryLocBuffer.WriteString("\n")
		loc.Key = c.ID + "_DEF"
		countryLocBuffer.WriteByte(' ')
		countryLocBuffer.WriteString(loc.Encode())
		countryLocBuffer.WriteString("\n")
		loc.Key = c.ID + "_ADJ"
		countryLocBuffer.WriteByte(' ')
		countryLocBuffer.WriteString(loc.Encode())
		countryLocBuffer.WriteString("\n\n")
	}
	err = util.WriteFileWithBOM(filepath.Join(modPath, "localisation", "simp_chinese", "tew_countries_auto_generate_l_simp_chinese.yml"), countryLocBuffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("生成国家名字文件成功！")

	fmt.Println("生成国家颜色文件中...")
	existColors := hashset.NewHashSetWith(stlslices.FlatMap(maps.Values(config.Countries), func(_ int, c *config.Country) [][3]uint8 {
		if c.Color.IsNone() {
			return nil
		}
		return [][3]uint8{{c.Color.MustValue()[0], c.Color.MustValue()[1], c.Color.MustValue()[2]}}
	})...)
	countryColors := make(map[string]*common.CountryColor, len(config.Countries))
	for _, c := range config.Countries {
		if c.Color.IsNone() {
			var cc color.Color
			for i := 0; i < 51; i++ {
				cc = colorful.Hsv(float64(rand.N(361)), float64(rand.N(71))/100, float64(50+rand.N(51))/100)
				r, g, b := util.GetRGB(cc)
				if !existColors.Contain([3]uint8{r, g, b}) {
					break
				}
				if i == 50 {
					return errors.New("can generate more color")
				}
			}
			r, g, b := util.GetRGB(cc)
			c.Color = optional.Some([3]uint8{r, g, b})
		}
		countryColors[c.ID] = &common.CountryColor{
			Country: c.ID,
			Color:   util.NewRGB(c.Color.MustValue()[0], c.Color.MustValue()[1], c.Color.MustValue()[2]),
			ColorUI: util.NewRGB(c.Color.MustValue()[0], c.Color.MustValue()[1], c.Color.MustValue()[2]),
		}
	}
	var countryColorBuffer bytes.Buffer
	for _, cc := range countryColors {
		countryColorBuffer.WriteString(cc.Encode())
		countryColorBuffer.WriteString("\n")
	}
	err = os.WriteFile(filepath.Join(modPath, "common", "countries", "colors.txt"), countryColorBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成国家颜色文件成功！")

	fmt.Println("生成不同国家类型颜色文件中...")
	cosmeticCountryColors := make(map[string]*common.CountryColor, len(config.Countries))
	countryTypes := map[string][3]uint8{
		"anarchism":    {255, 107, 0},
		"communism":    {255, 0, 0},
		"democratic":   {0, 0, 255},
		"conservatism": {0, 255, 255},
		"feudalism":    {192, 192, 192},
		"dictatorship": {255, 255, 0},
		"fascism":      {102, 51, 0},
	}
	for _, c := range config.Countries {
		for t, tc := range countryTypes {
			id := fmt.Sprintf("%s_type_%s", c.ID, t)
			cc := util.AlphaBlendColor(
				util.NewRGB(c.Color.MustValue()[0], c.Color.MustValue()[1], c.Color.MustValue()[2]),
				util.NewRGB(tc[0], tc[1], tc[2]),
				float32(rand.N(4)+4)/10,
			)
			// cc := color.GetRGBA{R: (c.Color.MustValue()[0] + tc[0]) / 2, G: (c.Color.MustValue()[1] + tc[1]) / 2, B: (c.Color.MustValue()[2] + tc[2]) / 2}
			cosmeticCountryColors[id] = &common.CountryColor{
				Country: id,
				Color:   cc,
				ColorUI: cc,
			}
		}
	}
	var cosmeticCountryColorBuffer bytes.Buffer
	for _, ccc := range cosmeticCountryColors {
		cosmeticCountryColorBuffer.WriteString(ccc.Encode())
		cosmeticCountryColorBuffer.WriteString("\n")
	}
	err = os.WriteFile(filepath.Join(modPath, "common", "countries", "cosmetic.txt"), cosmeticCountryColorBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成不同国家类型颜色文件成功！")

	fmt.Println("生成国家不同类型名字文件中...")
	countryTypeNameFormat := map[string]string{
		"anarchism":    "%s公社",
		"communism":    "%s社会主义共和国",
		"democratic":   "%s共和国",
		"conservatism": "%s王国",
		"feudalism":    "%s王国",
		"dictatorship": "%s国",
		"fascism":      "大%s帝国",
	}
	var cosmeticCountryLocBuffer bytes.Buffer
	cosmeticCountryLocBuffer.WriteString("l_simp_chinese:\n")
	for _, c := range config.Countries {
		for t, _ := range countryTypes {
			countryID := fmt.Sprintf("%s_type_%s", c.ID, t)
			loc := &localisation.Localisation{
				Key:   countryID,
				Index: optional.Some(0),
				Value: fmt.Sprintf(countryTypeNameFormat[t], c.Name),
			}
			cosmeticCountryLocBuffer.WriteByte(' ')
			cosmeticCountryLocBuffer.WriteString(loc.Encode())
			cosmeticCountryLocBuffer.WriteString("\n")
			loc.Key = countryID + "_DEF"
			cosmeticCountryLocBuffer.WriteByte(' ')
			cosmeticCountryLocBuffer.WriteString(loc.Encode())
			cosmeticCountryLocBuffer.WriteString("\n")
			loc.Key = countryID + "_ADJ"
			cosmeticCountryLocBuffer.WriteByte(' ')
			cosmeticCountryLocBuffer.WriteString(loc.Encode())
			cosmeticCountryLocBuffer.WriteString("\n")
		}
		cosmeticCountryLocBuffer.WriteString("\n")
	}
	err = util.WriteFileWithBOM(filepath.Join(modPath, "localisation", "simp_chinese", "tew_country_types_l_simp_chinese.yml.bak"), cosmeticCountryLocBuffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("生成国家不同类型名字文件成功！")

	fmt.Println("生成国家不同傀儡类型名字文件中...")
	autonomyTypeNameFormat := map[string]string{
		"dominion": "$OVERLORDADJ$属$NONIDEOLOGYADJ$自治领",
		"colony":   "$OVERLORDADJ$属$NONIDEOLOGYADJ$殖民政府",
		"puppet":   "$OVERLORDADJ$属$NONIDEOLOGYADJ$",
		"union":    "$OVERLORDADJ$-$NONIDEOLOGYADJ$邦",
		"division": "$OVERLORDADJ$-$NONIDEOLOGYADJ$军阀",
	}
	var autonomyCountryLocBuffer bytes.Buffer
	autonomyCountryLocBuffer.WriteString("l_simp_chinese:\n")
	for _, c := range config.Countries {
		for at, name := range autonomyTypeNameFormat {
			countryID := fmt.Sprintf("%s_tew_autonomy_%s", c.ID, at)
			loc := &localisation.Localisation{
				Key:   countryID,
				Index: optional.Some(0),
				Value: name,
			}
			autonomyCountryLocBuffer.WriteByte(' ')
			autonomyCountryLocBuffer.WriteString(loc.Encode())
			autonomyCountryLocBuffer.WriteString("\n")
			loc.Key = countryID + "_DEF"
			autonomyCountryLocBuffer.WriteByte(' ')
			autonomyCountryLocBuffer.WriteString(loc.Encode())
			autonomyCountryLocBuffer.WriteString("\n")
		}
		autonomyCountryLocBuffer.WriteString("\n")
	}
	err = util.WriteFileWithBOM(filepath.Join(modPath, "localisation", "simp_chinese", "tew_autonomy_name_l_simp_chinese copy.yml"), autonomyCountryLocBuffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("生成国家不同傀儡类型名字文件成功！")

	fmt.Println("生成国家动态变化脚本文件中...")
	var scriptedEffectBuffer bytes.Buffer
	var i int
	for _, c := range config.Countries {
		scriptedEffectBuffer.WriteString("\t\t")
		scriptedEffectBuffer.WriteString(stlbasic.Ternary(i == 0, "if", "else_if"))
		scriptedEffectBuffer.WriteString(" = {\n\t\t\tlimit = { original_tag = ")
		scriptedEffectBuffer.WriteString(c.ID)
		scriptedEffectBuffer.WriteString(" }\n")

		var j int
		for t, _ := range countryTypes {
			scriptedEffectBuffer.WriteString("\t\t\t")
			scriptedEffectBuffer.WriteString(stlbasic.Ternary(j == 0, "if", "else_if"))
			scriptedEffectBuffer.WriteString(" = {\n\t\t\t\tlimit = { has_country_flag = country_type_")
			scriptedEffectBuffer.WriteString(t)
			scriptedEffectBuffer.WriteString(" }\n\t\t\t\tset_cosmetic_tag = ")
			scriptedEffectBuffer.WriteString(c.ID)
			scriptedEffectBuffer.WriteString("_type_")
			scriptedEffectBuffer.WriteString(t)
			scriptedEffectBuffer.WriteString("\n\t\t\t}\n")
			j++
		}
		scriptedEffectBuffer.WriteString("\t\t}\n")
		i++
	}
	err = os.WriteFile("scripted_effects.txt", scriptedEffectBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成国家动态变化脚本文件成功！")
	return nil
}
