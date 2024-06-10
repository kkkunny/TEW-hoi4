package main

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/kkkunny/TEW-hoi4/sdk"
)

// type CountryHistory struct {
// 	Ideology string
// }
//
// func GetAllCountryHistories() (map[string]*CountryHistory, error) {
// 	dirPath := config.HOI4RootPath.Join("history", "countries")
// 	fileInfos, err := os.ReadDir(dirPath.String())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	histories, err := stlslices.FlatMapError(fileInfos, func(_ int, fileInfo os.DirEntry) ([]pair.Pair[string, *CountryHistory], error) {
// 		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".txt") {
// 			return nil, nil
// 		}
// 		countryTag := strings.TrimSpace(stlslices.First(strings.Split(fileInfo.Name(), "-")))
// 		fileData, err := os.ReadFile(dirPath.Join(fileInfo.Name()).String())
// 		if err != nil {
// 			return nil, err
// 		}
// 		partyMatchs := regexp.MustCompile(`ruling_party\s*=\s*(\S+)\s*`).FindStringSubmatch(string(fileData))
// 		if len(partyMatchs) < 2 {
// 			return nil, nil
// 		}
// 		return []pair.Pair[string, *CountryHistory]{pair.NewPair(countryTag, &CountryHistory{Ideology: partyMatchs[1]})}, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return stlslices.ToMap[pair.Pair[string, *CountryHistory], []pair.Pair[string, *CountryHistory], string, *CountryHistory, map[string]*CountryHistory](histories, func(p pair.Pair[string, *CountryHistory]) (string, *CountryHistory) {
// 		return p.First, p.Second
// 	}), nil
// }
//
// type CountryInfo struct {
// 	CountryTag
// 	CountryColor
// 	CountryLocalisation
// 	History optional.Optional[*CountryHistory]
// }
//
// func GetAllCountryInfos(ideologies []string) (map[string]*CountryInfo, error) {
// 	tags, err := GetAllCountryTags()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defs, err := GetAllCountryDefines(maps.Values(tags))
// 	if err != nil {
// 		return nil, err
// 	}
// 	colors, err := GetAllCountryColors()
// 	if err != nil {
// 		return nil, err
// 	}
// 	locs, err := GetAllCountryLocalisations(ideologies)
// 	if err != nil {
// 		return nil, err
// 	}
// 	histories, err := GetAllCountryHistories()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return stlmaps.MapError[string, *CountryTag, map[string]*CountryTag, string, *CountryInfo, map[string]*CountryInfo](tags, func(k string, tag *CountryTag) (string, *CountryInfo, error) {
// 		colorVal, ok := colors[k]
// 		if !ok {
// 			colorVal = &CountryColor{
// 				Color:   defs[k].Color,
// 				ColorUI: defs[k].Color,
// 			}
// 		}
// 		var history optional.Optional[*CountryHistory]
// 		if historyInfo, ok := histories[k]; ok {
// 			history = optional.Some(historyInfo)
// 		}
// 		return k, &CountryInfo{
// 			CountryTag:          *tag,
// 			CountryColor:        *colorVal,
// 			CountryLocalisation: locs[k],
// 			History:             history,
// 		}, nil
// 	})
// }
//
// type IdeologyInfo struct {
// 	ID    string
// 	Color color.RGBA
// }
//
// func GetAllIdeologies() (map[string]*IdeologyInfo, error) {
// 	ideologyData, err := os.ReadFile(config.HOI4RootPath.Join("common", "ideologies", "00_ideologies.txt").String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	ideologyStr := strings.ReplaceAll(string(ideologyData), "\n", " ")
//
// 	ideologyIDMatchs := regexp.MustCompile(`(\w+?)\s*=\s*\{\s*types\s*`).FindAllStringSubmatch(ideologyStr, -1)
// 	ideologyIDs, err := stlslices.MapError(ideologyIDMatchs, func(_ int, ideologyIDMatch []string) (string, error) {
// 		if len(ideologyIDMatch) < 2 {
// 			return "", errors.New("ideology format error")
// 		}
// 		return ideologyIDMatch[1], nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	ideologyColorMatchs := regexp.MustCompile(`color\s*=\s*(\w*?)\s*{\s*([\d.]+?)\s*([\d.]+?)\s*([\d.]+?)\s*}`).FindAllStringSubmatch(ideologyStr, -1)
// 	ideologyColors, err := stlslices.MapError(ideologyColorMatchs, func(_ int, ideologyColorMatch []string) (color.RGBA, error) {
// 		if len(ideologyColorMatch) < 4 {
// 			return stlbasic.Default[color.RGBA](), errors.New("ideology format error")
// 		}
// 		v1, _ := strconv.ParseFloat(ideologyColorMatch[2], 64)
// 		v2, _ := strconv.ParseFloat(ideologyColorMatch[3], 64)
// 		v3, _ := strconv.ParseFloat(ideologyColorMatch[4], 64)
// 		return util.Color(ideologyColorMatch[1], v1, v2, v3)
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(ideologyIDs) != len(ideologyColors) {
// 		return nil, errors.New("ideology format error")
// 	}
// 	infos := stlslices.Map(ideologyIDs, func(i int, ideologyID string) *IdeologyInfo {
// 		return &IdeologyInfo{
// 			ID:    ideologyID,
// 			Color: ideologyColors[i],
// 		}
// 	})
// 	return stlslices.ToMap[*IdeologyInfo, []*IdeologyInfo, string, *IdeologyInfo, map[string]*IdeologyInfo](infos, func(info *IdeologyInfo) (string, *IdeologyInfo) {
// 		return info.ID, info
// 	}), nil
// }
//
// type CountryIdeologyLocalisation struct {
// 	Name string
// 	DEF  string
// 	ADJ  string
// }
//
// type CountryLocalisation map[string]*CountryIdeologyLocalisation
//
// func GetAllCountryLocalisations(ideologies []string) (map[string]CountryLocalisation, error) {
// 	res := make(map[string]CountryLocalisation)
// 	locDir := config.HOI4RootPath.Join("localisation", "simp_chinese")
// 	locFileInfos, err := os.ReadDir(locDir.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, locFileInfo := range locFileInfos {
// 		if locFileInfo.IsDir() || !strings.HasSuffix(locFileInfo.Name(), ".yml") {
// 			continue
// 		}
// 		locData, err := os.ReadFile(locDir.Join(locFileInfo.Name()).String())
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		ideologyNames := stlslices.Map(ideologies, func(_ int, e string) string {
// 			return "_" + e
// 		})
// 		re := fmt.Sprintf(`\s+([A-Z0-9]{3})(%s)?(_DEF|_ADJ)?\s*:\s*"(.*?)"`, strings.Join(ideologyNames, "|"))
// 		locMatchs := regexp.MustCompile(re).FindAllStringSubmatch(strings.ReplaceAll(string(locData), "\n", " "), -1)
// 		for _, locMatch := range locMatchs {
// 			if len(locMatch) < 5 {
// 				continue
// 			}
// 			countryLoc, ok := res[locMatch[1]]
// 			if !ok {
// 				res[locMatch[1]] = make(CountryLocalisation)
// 				countryLoc = res[locMatch[1]]
// 			}
// 			ideologyLoc, ok := countryLoc[strings.TrimPrefix(locMatch[2], "_")]
// 			if !ok {
// 				countryLoc[strings.TrimPrefix(locMatch[2], "_")] = new(CountryIdeologyLocalisation)
// 				ideologyLoc = countryLoc[strings.TrimPrefix(locMatch[2], "_")]
// 			}
// 			switch strings.TrimPrefix(locMatch[3], "_") {
// 			case "":
// 				ideologyLoc.Name = locMatch[4]
// 			case "DEF":
// 				ideologyLoc.DEF = locMatch[4]
// 			case "ADJ":
// 				ideologyLoc.ADJ = locMatch[4]
// 			}
// 		}
// 	}
// 	return res, nil
// }
//
// func WriteCosmeticTagDefine(prefix string, modPath stlos.FilePath, countryInfos []*CountryInfo, ideologyInfos []*IdeologyInfo) error {
// 	file, err := os.OpenFile(modPath.Join("common", "countries", "cosmetic.txt").String(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	for i, countryInfo := range countryInfos {
// 		for j, ideologyInfo := range ideologyInfos {
// 			r, g, b := countryInfo.Color.R, countryInfo.Color.G, countryInfo.Color.B
// 			uir, uig, uib := countryInfo.ColorUI.R, countryInfo.ColorUI.G, countryInfo.ColorUI.B
// 			if countryInfo.History.IsNone() || countryInfo.History.MustValue().Ideology != ideologyInfo.ID {
// 				var h, s, v float64
// 				switch ideologyInfo.ID {
// 				case "communism":
// 					h = float64(util.RandomInt(0, 41) - 20)
// 					if h < 0 {
// 						h += 360
// 					}
// 					s = float64(util.RandomInt(80, 101)) / 100
// 					v = float64(util.RandomInt(80, 101)) / 100
// 				case "fascism":
// 					h = float64(util.RandomInt(25, 51))
// 					s = float64(util.RandomInt(60, 101)) / 100
// 					v = float64(util.RandomInt(70, 91)) / 100
// 				case "neutrality":
// 					h = float64(util.RandomInt(0, 361))
// 					s = float64(util.RandomInt(40, 51)) / 100
// 					v = float64(util.RandomInt(50, 71)) / 100
// 				case "democratic":
// 					h = float64(util.RandomInt(200, 251))
// 					s = float64(util.RandomInt(70, 91)) / 100
// 					v = float64(util.RandomInt(80, 91)) / 100
// 				}
// 				r, g, b = util.HSV2RGB(h, s, v)
// 				uir, uig, uib = r, g, b
// 			}
//
// 			_, err = fmt.Fprintf(
// 				file,
// 				"%s_%s_%s = {\n\tcolor = rgb { %d %d %d }\n\tcolor_ui = rgb { %d %d %d }\n}",
// 				prefix,
// 				countryInfo.Tag,
// 				ideologyInfo.ID,
// 				r,
// 				g,
// 				b,
// 				uir,
// 				uig,
// 				uib,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			if j != len(ideologyInfos)-1 {
// 				_, err = fmt.Fprintf(file, "\n")
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		if i != len(countryInfos)-1 {
// 			_, err = fmt.Fprintf(file, "\n\n")
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
//
// func WriteCosmeticTagScriptedEffect(prefix string, modPath stlos.FilePath, countryInfos []*CountryInfo, ideologyInfos []*IdeologyInfo) error {
// 	file, err := os.OpenFile(modPath.Join("common", "scripted_effects", fmt.Sprintf("00_%s_scripted_effects.txt", prefix)).String(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	_, err = fmt.Fprintf(file, "%s_change_country_color_by_ideology = {\n", prefix)
// 	if err != nil {
// 		return err
// 	}
// 	for i, countryInfo := range countryInfos {
// 		_, err = fmt.Fprintf(file, "\t%s = {\n\t\tlimit = { original_tag = %s }\n", stlbasic.Ternary(i == 0, "if", "else_if"), countryInfo.Tag)
// 		if err != nil {
// 			return err
// 		}
// 		for j, ideologyInfo := range ideologyInfos {
// 			_, err = fmt.Fprintf(
// 				file,
// 				"\t\t%s = {\n\t\t\tlimit = {\n\t\t\t\tOR = {\n\t\t\t\t\thas_government = %s\n\t\t\t\t}\n\t\t\t}\n\t\t\tset_cosmetic_tag = %s_%s_%s\n\t\t}",
// 				stlbasic.Ternary(j == 0, "if", "else_if"),
// 				ideologyInfo.ID,
// 				prefix,
// 				countryInfo.Tag,
// 				ideologyInfo.ID,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			if j != len(ideologyInfos)-1 {
// 				_, err = fmt.Fprintf(file, "\n")
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		_, err = fmt.Fprintf(file, "\n\t}")
// 		if err != nil {
// 			return err
// 		}
// 		if i != len(countryInfos)-1 {
// 			_, err = fmt.Fprintf(file, "\n")
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	_, err = fmt.Fprintf(file, "\n}")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func WriteCosmeticTagLocalisation(prefix string, modPath stlos.FilePath, countryInfos []*CountryInfo, ideologyInfos []*IdeologyInfo) error {
// 	file, err := os.OpenFile(modPath.Join("localisation", "simp_chinese", fmt.Sprintf("%s_countries_cosmetic_l_simp_chinese.yml", prefix)).String(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	_, err = fmt.Fprintf(file, "l_simp_chinese:\n")
// 	if err != nil {
// 		return err
// 	}
// 	for i, countryInfo := range countryInfos {
// 		for j, ideologyInfo := range ideologyInfos {
// 			ideoLogyLoc := stlbasic.TernaryAction(slices.Contains(maps.Keys(countryInfo.CountryLocalisation), ideologyInfo.ID), func() *CountryIdeologyLocalisation {
// 				return countryInfo.CountryLocalisation[ideologyInfo.ID]
// 			}, func() *CountryIdeologyLocalisation {
// 				if stlmaps.Empty(countryInfo.CountryLocalisation) {
// 					return &CountryIdeologyLocalisation{}
// 				}
// 				v, ok := countryInfo.CountryLocalisation[""]
// 				if ok {
// 					return v
// 				}
// 				_, v = stlmaps.Random(countryInfo.CountryLocalisation)
// 				return v
// 			})
// 			_, err = fmt.Fprintf(
// 				file,
// 				" %s_%s_%s:0 \"%s\"\n %s_%s_%s_DEF:0 \"%s\"\n %s_%s_%s_ADJ:0 \"%s\"",
// 				prefix,
// 				countryInfo.Tag,
// 				ideologyInfo.ID,
// 				ideoLogyLoc.Name,
// 				prefix,
// 				countryInfo.Tag,
// 				ideologyInfo.ID,
// 				ideoLogyLoc.DEF,
// 				prefix,
// 				countryInfo.Tag,
// 				ideologyInfo.ID,
// 				ideoLogyLoc.ADJ,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			if j != len(ideologyInfos)-1 {
// 				_, err = fmt.Fprintf(file, "\n")
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		if i != len(countryInfos)-1 {
// 			_, err = fmt.Fprintf(file, "\n")
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
//
// func WriteCosmeticTagFlag(prefix string, modPath stlos.FilePath, countryInfos []*CountryInfo, ideologyInfos []*IdeologyInfo) error {
// 	flagDir := config.HOI4RootPath.Join("gfx", "flags")
// 	modFlagDir := modPath.Join("gfx", "flags")
// 	suffix := append([]string{""}, stlslices.Map(ideologyInfos, func(_ int, ideology *IdeologyInfo) string {
// 		return "_" + ideology.ID
// 	})...)
//
// 	for _, countryInfo := range countryInfos {
// 		for _, s := range suffix {
// 			flagPath := flagDir.Join(fmt.Sprintf("%s%s.tga", countryInfo.Tag, s))
// 			exist, err := stlos.Exist(flagPath)
// 			if err != nil {
// 				return err
// 			}
// 			if !exist {
// 				continue
// 			}
// 			err = util.CopyFile(flagPath.String(), modFlagDir.Join(fmt.Sprintf("%s.tga", countryInfo.Tag)).String())
// 			if err != nil {
// 				return err
// 			}
// 			err = util.ResizeAndCopyTgaImage(flagPath.String(), modFlagDir.Join("medium", fmt.Sprintf("%s.tga", countryInfo.Tag)).String(), 41, 26)
// 			if err != nil {
// 				return err
// 			}
// 			err = util.ResizeAndCopyTgaImage(flagPath.String(), modFlagDir.Join("small", fmt.Sprintf("%s.tga", countryInfo.Tag)).String(), 10, 7)
// 			if err != nil {
// 				return err
// 			}
// 			break
// 		}
// 		for _, ideologyInfo := range ideologyInfos {
// 			for _, s := range []string{fmt.Sprintf("_%s", ideologyInfo.ID), ""} {
// 				flagPath := flagDir.Join(fmt.Sprintf("%s%s.tga", countryInfo.Tag, s))
// 				exist, err := stlos.Exist(flagPath)
// 				if err != nil {
// 					return err
// 				}
// 				if !exist {
// 					continue
// 				}
// 				err = util.CopyFile(flagPath.String(), modFlagDir.Join(fmt.Sprintf("%s_%s_%s.tga", prefix, countryInfo.Tag, ideologyInfo.ID)).String())
// 				if err != nil {
// 					return err
// 				}
// 				err = util.ResizeAndCopyTgaImage(flagPath.String(), modFlagDir.Join("medium", fmt.Sprintf("%s_%s_%s.tga", prefix, countryInfo.Tag, ideologyInfo.ID)).String(), 41, 26)
// 				if err != nil {
// 					return err
// 				}
// 				err = util.ResizeAndCopyTgaImage(flagPath.String(), modFlagDir.Join("small", fmt.Sprintf("%s_%s_%s.tga", prefix, countryInfo.Tag, ideologyInfo.ID)).String(), 10, 7)
// 				if err != nil {
// 					return err
// 				}
// 				break
// 			}
// 		}
// 	}
// 	return nil
// }

func main() {
	// modPath := stlos.NewFilePath("mod/DynamicCountryColor")
	// ideologies, err := GetAllIdeologies()
	// if err != nil {
	// 	panic(err)
	// }
	// countriInfos, err := GetAllCountryInfos(maps.Keys(ideologies))
	// if err != nil {
	// 	panic(err)
	// }
	// err = WriteCosmeticTagDefine("DCC", modPath, maps.Values(countriInfos), maps.Values(ideologies))
	// if err != nil {
	// 	panic(err)
	// }
	// err = WriteCosmeticTagDefine("DCC", modPath, maps.Values(countriInfos), maps.Values(ideologies))
	// if err != nil {
	// 	panic(err)
	// }

	// // stateContinentMap, err := _map.ParseStateDef(filepath.Join(modPath, "map", "definition.csv"))
	// stateInfos, err := os.ReadDir(filepath.Join(config.TEWRootPath, "history", "states"))
	// if err != nil {
	// 	panic(err)
	// }
	// for _, stateInfo := range stateInfos {
	// 	if stateInfo.IsDir() {
	// 		continue
	// 	}
	// 	fp := filepath.Join(config.TEWRootPath, "history", "states", stateInfo.Name())
	// 	state, err := history.ParseState(fp)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	// if stateDef, ok := stateContinentMap[state.ID]; ok && stateDef.ContinentID == 4 {
	// 	// 	cores := hashset.NewHashSetWith(state.History.Cores...)
	// 	// 	cores.Add("NAM")
	// 	// 	state.History.Cores = cores.ToSlice().ToSlice()
	// 	// } else {
	// 	// 	continue
	// 	// }
	// 	switch {
	// 	// 删除
	// 	// case stlslices.Contain(state.History.Cores, "GAL") || stlslices.Contain(state.History.Claims, "GAL"):
	// 	// 	claims := hashset.NewHashSetWith(state.History.Claims...)
	// 	// 	claims.Remove("GAL")
	// 	// 	state.History.Claims = claims.ToSlice().ToSlice()
	// 	// 	cores := hashset.NewHashSetWith(state.History.Cores...)
	// 	// 	cores.Remove("GAL")
	// 	// 	state.History.Cores = cores.ToSlice().ToSlice()
	// 	case stlslices.ContainAny(state.History.Cores, "WLS"):
	// 		cores := hashset.NewHashSetWith(state.History.Cores...)
	// 		cores.Add("WRM")
	// 		cores.Add("RME")
	// 		state.History.Cores = cores.ToSlice().ToSlice()
	// 		// claims := hashset.NewHashSetWith(state.History.Claims...)
	// 		// claims.Add("TKS")
	// 		// state.History.Claims = claims.ToSlice().ToSlice()
	// 	default:
	// 		continue
	// 	}
	// 	fmt.Println(stateInfo.Name())
	//
	// 	err = os.WriteFile(fp, []byte(state.Encode()), 0666)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return sdk.RefreshCountries()
	})
	// eg.Go(func() error {
	// 	flagPath := filepath.Join(config.TEWRootPath, "gfx", "flags")
	// 	flagInfos, err := os.ReadDir(flagPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	smallFlagPath := filepath.Join(flagPath, "small")
	// 	err = os.RemoveAll(smallFlagPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = os.MkdirAll(smallFlagPath, 0666)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, flagInfo := range flagInfos {
	// 		if flagInfo.IsDir() {
	// 			continue
	// 		}
	// 		fp := filepath.Join(flagPath, flagInfo.Name())
	// 		err = util.ResizeAndCopyTgaImage(fp, filepath.Join(smallFlagPath, flagInfo.Name()), 10, 7)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })
	// eg.Go(func() error {
	// 	flagPath := filepath.Join(config.TEWRootPath, "gfx", "flags")
	// 	flagInfos, err := os.ReadDir(flagPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	mediumFlagPath := filepath.Join(flagPath, "medium")
	// 	err = os.RemoveAll(mediumFlagPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = os.MkdirAll(mediumFlagPath, 0666)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, flagInfo := range flagInfos {
	// 		if flagInfo.IsDir() {
	// 			continue
	// 		}
	// 		fp := filepath.Join(flagPath, flagInfo.Name())
	// 		err = util.ResizeAndCopyTgaImage(fp, filepath.Join(mediumFlagPath, flagInfo.Name()), 41, 26)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
