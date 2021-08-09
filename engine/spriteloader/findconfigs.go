package spriteloader

import (
	"path/filepath"
	"os"
	"strings"
)

func findConfigPaths() []string {
	paths := []string{}
	filepath.Walk(
		"./assets/sprite",
		func(path string, info os.FileInfo, er error) error {
			if strings.HasSuffix(path, ".yaml") {
				paths = append(paths, path)
			}
			return nil
		},
	)
	return paths
}

func readConfigs(paths []string) [][]SpriteID {
	configs := make([][]SpriteID,len(paths))
	for idx,path := range paths {
		configs[idx] = RegisterSpritesFromConfig(path)
	}
	return configs
}

func LoadSpritesFromConfigs() {
	configPaths := findConfigPaths()
	readConfigs(configPaths)
}
