package config

import "path/filepath"

var (
	HOI4RootPath  = `F:\SteamLibrary\steamapps\common\Hearts of Iron IV`
	HOI4ModPath   = `C:\Users\14012\Documents\Paradox Interactive\Hearts of Iron IV\mod`
	HOI4MyModPath = `W:\code\go\github.com\kkkunny\TEW-hoi4\mod`
	TEWRootPath   = filepath.Join(HOI4MyModPath, "TheEmptyWorld")
)
