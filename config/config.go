package config

import (
	"path/filepath"
	"os"
	"strings"
)

// ext should look like ".yaml" or ".png" - dot must be included
func FindPathsWithExt(root,ext string) []string {
	paths := []string{}
	filepath.Walk(
		root,
		func(path string, info os.FileInfo, er error) error {
			if strings.HasSuffix(path, ext) {
				paths = append(paths, path)
			}
			return nil
		},
	)
	return paths
}
