package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	stlbasic "github.com/kkkunny/stl/basic"
	"github.com/kkkunny/stl/container/hashset"
	"github.com/kkkunny/stl/container/linkedhashmap"
	stlmaps "github.com/kkkunny/stl/container/maps"
	"github.com/kkkunny/stl/container/optional"
	"github.com/kkkunny/stl/container/pair"
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
	countryTypes := linkedhashmap.NewLinkedHashMapWith[string, [3]uint8](
		"anarchism", [3]uint8{255, 107, 0},
		"communism", [3]uint8{255, 0, 0},
		"democratic", [3]uint8{0, 0, 255},
		"conservatism", [3]uint8{0, 255, 255},
		"feudalism", [3]uint8{192, 192, 192},
		"dictatorship", [3]uint8{255, 255, 0},
		"fascism", [3]uint8{102, 51, 0},
	)
	for _, c := range config.Countries {
		for iter := countryTypes.Iterator(); iter.Next(); {
			tc := iter.Value().Second
			id := fmt.Sprintf("%s_type_%s", c.ID, iter.Value().First)
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

	locs, err := localisation.ParseChineseLocalisationDir(config.TEWRootPath)
	if err != nil {
		return err
	}

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
		for iter := countryTypes.Iterator(); iter.Next(); {
			countryID := fmt.Sprintf("%s_type_%s", c.ID, iter.Value().First)
			loc := &localisation.Localisation{
				Key:   countryID,
				Index: optional.Some(0),
				Value: stlbasic.Ternary(stlmaps.ContainKey(locs, countryID), locs[countryID].Value, fmt.Sprintf(countryTypeNameFormat[iter.Value().First], c.Name)),
			}
			cosmeticCountryLocBuffer.WriteByte(' ')
			cosmeticCountryLocBuffer.WriteString(loc.Encode())
			cosmeticCountryLocBuffer.WriteString("\n")
			loc.Key = countryID + "_DEF"
			cosmeticCountryLocBuffer.WriteByte(' ')
			cosmeticCountryLocBuffer.WriteString(loc.Encode())
			cosmeticCountryLocBuffer.WriteString("\n")
		}
		cosmeticCountryLocBuffer.WriteString("\n")
	}
	err = util.WriteFileWithBOM(filepath.Join(modPath, "localisation", "simp_chinese", "tew_country_types_auto_generate_l_simp_chinese.yml"), cosmeticCountryLocBuffer.Bytes())
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

	var findSonCountries func(tag string) []*config.Country
	findSonCountries = func(tag string) []*config.Country {
		c, ok := config.Countries[tag]
		if !ok || config.Countries[tag].Sons.IsNone() || stlslices.Empty(config.Countries[tag].Sons.MustValue()) {
			return nil
		}
		return stlslices.RemoveRepeat(stlslices.Filter(stlslices.FlatMap(c.Sons.MustValue(), func(_ int, e string) []*config.Country {
			return append(findSonCountries(e), config.Countries[e])
		}), func(_ int, e *config.Country) bool {
			return e != nil
		}))
	}

	canUpgradedCountries := stlslices.ToMap(stlslices.Filter(maps.Values(config.Countries), func(_ int, c *config.Country) bool {
		return c.Sons.IsSome() && !stlslices.Empty(c.Sons.MustValue())
	}), func(e *config.Country) (string, *config.Country) {
		return e.ID, e
	})
	sonTag2ParentTags := stlslices.FlatMap(maps.Values(canUpgradedCountries), func(_ int, c *config.Country) []pair.Pair[*config.Country, *config.Country] {
		return stlslices.Filter(stlslices.FlatMap(c.Sons.MustValue(), func(_ int, sonTag string) []pair.Pair[*config.Country, *config.Country] {
			return stlslices.Map(append(findSonCountries(sonTag), config.Countries[sonTag]), func(_ int, son *config.Country) pair.Pair[*config.Country, *config.Country] {
				return pair.NewPair(son, c)
			})
		}), func(_ int, e pair.Pair[*config.Country, *config.Country]) bool {
			return e.First != nil
		})
	})
	sonTag2ParentTagsMap := make(map[string]*hashset.HashSet[string], len(sonTag2ParentTags))
	for _, p := range sonTag2ParentTags {
		if sonTag2ParentTagsMap[p.First.ID] == nil {
			sonTag2ParentTagsMap[p.First.ID] = stlbasic.Ptr(hashset.NewHashSet[string]())
		}
		sonTag2ParentTagsMap[p.First.ID].Add(p.Second.ID)
	}

	fmt.Println("生成国家成立脚本文件中...")
	var tagScriptedEffectBuffer bytes.Buffer
	tagScriptedEffectBuffer.WriteString("ideas = {\n\tcountry_tag = {\n\t\tlaw = yes\n\n\t\tcountry_tag_default = {\n\t\t\ton_add = {\n\t\t\t\ttew_update_country_type = yes\n\t\t\t}\n\n\t\t\tai_will_do = { factor = 0 }\n\n\t\t\tcancel_if_invalid = yes\n\t\t\tdefault = yes\n\t\t}")
	for _, c := range canUpgradedCountries {
		tagScriptedEffectBuffer.WriteString("\n\n\t\t")
		tagScriptedEffectBuffer.WriteString("country_tag_")
		tagScriptedEffectBuffer.WriteString(c.ID)
		tagScriptedEffectBuffer.WriteString(" = {\n\t\t\tallowed = {\n\t\t\t\tOR = {\n\t\t\t\t\t")
		tagScriptedEffectBuffer.WriteString(strings.Join(stlslices.Map(findSonCountries(c.ID), func(_ int, e *config.Country) string { return "original_tag = " + e.ID }), "\n\t\t\t\t\t"))
		tagScriptedEffectBuffer.WriteString(fmt.Sprintf("\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tavailable = {\n\t\t\t\ttew_can_ndependent_diplomacy = yes\n\t\t\t\tNOT = {\n\t\t\t\t\tcountry_exists = %s\n\t\t\t\t\tany_other_country = {\n\t\t\t\t\t\tlimit = { exists = yes }\n\t\t\t\t\t\thas_idea = country_tag_%s\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\t%s = {\n\t\t\t\t\tset_temp_variable = { tew_than_number = %d }\n\t\t\t\t\ttew_self_or_puppet_owns_gte = yes\n\t\t\t\t}\n\t\t\t}\n\n\t\t\ton_add = {\n\t\t\t\ttew_update_country_type = yes\n\t\t\t}\n\n\t\t\tai_will_do = {\n\t\t\t\tfactor = 100", c.ID, c.ID, c.ID, c.UpgradeRatio.ValueWith(70)))
		sonTags := stlslices.Map(findSonCountries(c.ID), func(_ int, e *config.Country) string { return e.ID })
		conflictCountries := hashset.NewHashSet[string]()
		for _, scTag := range sonTags {
			parents, ok := sonTag2ParentTagsMap[scTag]
			if !ok {
				continue
			}
			for iter := parents.Iterator(); iter.Next(); {
				if iter.Value() == c.ID {
					continue
				} else if stlslices.Contain(sonTags, iter.Value()) {
					continue
				} else if conflictCountries.Contain(iter.Value()) {
					continue
				}
				conflictCountries.Add(iter.Value())
				tagScriptedEffectBuffer.WriteString(fmt.Sprintf("\n\t\t\t\tmodifier = {\n\t\t\t\t\tfactor = 0\n\t\t\t\t\thas_idea = country_tag_%s\n\t\t\t\t}", iter.Value()))
			}
		}
		tagScriptedEffectBuffer.WriteString("\n\t\t\t}\n\n\t\t\tcancel_if_invalid = yes\n\t\t}")
	}
	tagScriptedEffectBuffer.WriteString("\n}")
	err = os.WriteFile(filepath.Join(modPath, "common", "ideas", "tew_attr_country_tag_auto_generate.txt"), tagScriptedEffectBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成国家成立脚本文件成功！")

	fmt.Println("生成国家动态变化脚本文件中...")
	var scriptedEffectBuffer bytes.Buffer
	scriptedEffectBuffer.WriteString("# 更新国家类型\n#param: THIS\ntew_update_country_type = {\n\t# clear flag\n\tclr_country_flag = country_type_none\n\tclr_country_flag = country_type_anarchism\n\tclr_country_flag = country_type_communism\n\tclr_country_flag = country_type_democratic\n\tclr_country_flag = country_type_conservatism\n\tclr_country_flag = country_type_feudalism\n\tclr_country_flag = country_type_dictatorship\n\tclr_country_flag = country_type_fascism\n\n\t# set flag\n\tif = {\n\t\tlimit = { tew_can_ndependent_diplomacy = no }\n\t\tset_country_flag = country_type_none\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\thas_idea = gov_anarchist_commune\n\t\t}\n\t\tset_country_flag = country_type_anarchism\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\tOR = {\n\t\t\t\thas_idea = gov_communist_dictatorship\n\t\t\t\thas_idea = gov_communist_republic\n\t\t\t}\n\t\t}\n\t\tset_country_flag = country_type_communism\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\tOR = {\n\t\t\t\thas_idea = gov_presidential_republic\n\t\t\t\thas_idea = gov_parliamentary_republic\n\t\t\t\thas_idea = gov_committee_republic\n\t\t\t}\n\t\t}\n\t\tset_country_flag = country_type_democratic\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\thas_idea = gov_parliamentary_constitutional_monarchy\n\t\t}\n\t\tset_country_flag = country_type_conservatism\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\tOR = {\n\t\t\t\thas_idea = gov_dualist_constitutional_monarchy\n\t\t\t\thas_idea = gov_absolute_monarchy\n\t\t\t}\n\t\t}\n\t\tset_country_flag = country_type_feudalism\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\tOR = {\n\t\t\t\thas_idea = gov_presidential_dictatorship\n\t\t\t\thas_idea = gov_parliamentary_dictatorship\n\t\t\t\thas_idea = gov_military_dictatorship\n\t\t\t}\n\t\t}\n\t\tset_country_flag = country_type_dictatorship\n\t}\n\telse_if = {\n\t\tlimit = {\n\t\t\tOR = {\n\t\t\t\thas_idea = gov_fascist_republic\n\t\t\t\thas_idea = gov_fascist_dictatorship\n\t\t\t}\n\t\t}\n\t\tset_country_flag = country_type_fascism\n\t}\n\n\t# set cosmetic tag\n\tif = {\n\t\tlimit = { has_country_flag = country_type_none }\n\t\tdrop_cosmetic_tag = yes\n\t}\n\telse_if = {\n\t\tlimit = { NOT = { has_idea = country_tag_default } }\n")
	var i int
	for _, c := range canUpgradedCountries {
		scriptedEffectBuffer.WriteString("\t\t")
		scriptedEffectBuffer.WriteString(stlbasic.Ternary(i == 0, "if", "else_if"))
		scriptedEffectBuffer.WriteString(" = {\n\t\t\tlimit = { has_idea = country_tag_")
		scriptedEffectBuffer.WriteString(c.ID)
		scriptedEffectBuffer.WriteString(" }\n")

		var j int
		for iter := countryTypes.Iterator(); iter.Next(); {
			scriptedEffectBuffer.WriteString("\t\t\t")
			scriptedEffectBuffer.WriteString(stlbasic.Ternary(j == 0, "if", "else_if"))
			scriptedEffectBuffer.WriteString(" = {\n\t\t\t\tlimit = { has_country_flag = country_type_")
			scriptedEffectBuffer.WriteString(iter.Value().First)
			scriptedEffectBuffer.WriteString(" }\n\t\t\t\tset_cosmetic_tag = ")
			scriptedEffectBuffer.WriteString(c.ID)
			scriptedEffectBuffer.WriteString("_type_")
			scriptedEffectBuffer.WriteString(iter.Value().First)
			scriptedEffectBuffer.WriteString("\n\t\t\t}\n")
			j++
		}
		scriptedEffectBuffer.WriteString("\t\t}\n")
		i++
	}
	scriptedEffectBuffer.WriteString("\t}\n\telse = {\n")
	i = 0
	for _, c := range config.Countries {
		scriptedEffectBuffer.WriteString("\t\t")
		scriptedEffectBuffer.WriteString(stlbasic.Ternary(i == 0, "if", "else_if"))
		scriptedEffectBuffer.WriteString(" = {\n\t\t\tlimit = { original_tag = ")
		scriptedEffectBuffer.WriteString(c.ID)
		scriptedEffectBuffer.WriteString(" }\n")

		var j int
		for iter := countryTypes.Iterator(); iter.Next(); {
			scriptedEffectBuffer.WriteString("\t\t\t")
			scriptedEffectBuffer.WriteString(stlbasic.Ternary(j == 0, "if", "else_if"))
			scriptedEffectBuffer.WriteString(" = {\n\t\t\t\tlimit = { has_country_flag = country_type_")
			scriptedEffectBuffer.WriteString(iter.Value().First)
			scriptedEffectBuffer.WriteString(" }\n\t\t\t\tset_cosmetic_tag = ")
			scriptedEffectBuffer.WriteString(c.ID)
			scriptedEffectBuffer.WriteString("_type_")
			scriptedEffectBuffer.WriteString(iter.Value().First)
			scriptedEffectBuffer.WriteString("\n\t\t\t}\n")
			j++
		}
		scriptedEffectBuffer.WriteString("\t\t}\n")
		i++
	}
	scriptedEffectBuffer.WriteString("\t}\n}\n")
	err = os.WriteFile(filepath.Join(modPath, "common", "scripted_effects", "tew_tag_scripted_effects_auto_generate.txt"), scriptedEffectBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成国家动态变化脚本文件成功！")

	fmt.Println("生成可变身国家名字文件中...")
	var canUpgradeCountryLocBuffer bytes.Buffer
	canUpgradeCountryLocBuffer.WriteString("l_simp_chinese:\n country_tag:0 \"国家\"\n idea_group_country_tag:0 \"国家\"\n idea_group_country_tag_desc:0 \"国家\"\n country_tag_default:0 \"默认\"\n")
	for _, c := range canUpgradedCountries {
		canUpgradeCountryLocBuffer.WriteString(" country_tag_")
		canUpgradeCountryLocBuffer.WriteString(c.ID)
		canUpgradeCountryLocBuffer.WriteString(":0 \"")
		canUpgradeCountryLocBuffer.WriteString(c.Name)
		canUpgradeCountryLocBuffer.WriteString("\"\n")
	}
	err = util.WriteFileWithBOM(filepath.Join(modPath, "localisation", "simp_chinese", "tew_country_tag_auto_generate_l_simp_chinese.yml"), canUpgradeCountryLocBuffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("生成可变身国家名字文件成功！")

	fmt.Println("生成可变身国家图标文件中...")
	var canUpgradeCountryIconBuffer bytes.Buffer
	canUpgradeCountryIconBuffer.WriteString("spriteTypes = {")
	for _, c := range canUpgradedCountries {
		canUpgradeCountryIconBuffer.WriteString(fmt.Sprintf("\n\tspriteType = {\n\t\tname = \"GFX_idea_country_tag_%s\"\n\t\ttexturefile = \"gfx\\\\flags\\\\medium\\\\%s.tga\"\n\t}", c.ID, c.ID))
		canUpgradeCountryIconBuffer.WriteString("\n")
	}
	canUpgradeCountryIconBuffer.WriteString("}")
	err = os.WriteFile(filepath.Join(modPath, "interface", "tew_country_tga_auto_generate.gfx"), canUpgradeCountryIconBuffer.Bytes(), 0666)
	if err != nil {
		return err
	}
	fmt.Println("生成可变身国家图标文件成功！")
	return nil
}
