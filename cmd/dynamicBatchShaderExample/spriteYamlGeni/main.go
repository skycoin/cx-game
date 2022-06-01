package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/Texture"
)

var imagesPath = "../../../assets/sprite/spineboy"
var names []string

func main() {
	MakeYamlOfAllSpritesInPath(imagesPath)
}

func MakeYamlOfAllSpritesInPath(path string) {
	paths := findConfigPaths(path)

	for i, path := range paths {
		tex := Texture.ReadTextureData(path)
		tex.M_name = names[i]
		MakeYaml(*tex)
		//	break
	}

	fmt.Println("Done make Yaml Files")
}

type SpriteConfig struct {
	Name       string `yaml:"name"`
	Width      int    `yaml:"width"`
	Height     int    `yaml:"height"`
	CellWidth  int    `yaml:"cellwidth"`
	CellHeight int    `yaml:"cellheight"`
	Autoname   string `yaml:"autoname"`
}

func MakeYaml(sprite Texture.Texture) {
	fmt.Println(sprite)

	SpriteData := SpriteConfig{Name: sprite.M_name,
		Width:      sprite.M_width,
		Height:     sprite.M_height,
		CellWidth:  sprite.M_width,
		CellHeight: sprite.M_height,
		Autoname:   "index",
	}

	data, err := yaml.Marshal(SpriteData)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(SpriteData.Name)
	err2 := ioutil.WriteFile(imagesPath+"/"+sprite.M_name+".yaml", data, 0)

	if err2 != nil {

		log.Fatal(err2)
	}

	fmt.Println("data written")
}

func findConfigPaths(spritesPath string) []string {
	paths := []string{}

	filepath.Walk(
		spritesPath,
		func(path string, info os.FileInfo, er error) error {
			fmt.Println(strings.TrimSuffix(info.Name(), ".png"))
			if strings.HasSuffix(path, ".png") {
				names = append(names, strings.TrimSuffix(info.Name(), ".png"))
				paths = append(paths, path)
			}
			return nil
		},
	)

	return paths
}
