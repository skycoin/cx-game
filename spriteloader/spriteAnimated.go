package spriteloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

/*
	this is base struct for sprite animated
*/
type Frame struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type SpriteSourceSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type SourceSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Frames struct {
	Name             string
	Action           string
	Order            int
	Frame            Frame            `json:"frames"`
	Rotated          bool             `json:"rotated"`
	Trimmed          bool             `json:"trimmed"`
	SpriteSourceSize SpriteSourceSize `json:"spriteSourceSize"`
	SourceSize       SourceSize       `json:"sourceSize"`
	Duration         int              `json:"duration"`
}

type MetaSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type FrameTagItem struct {
	Name      string `json:"name"`
	From      int    `json:"from"`
	To        int    `json:"to"`
	Direction string `json:"direction"`
	Color     string `json:"color"`
}

type Meta struct {
	Image     string         `json:"image"`
	Format    string         `json:"format"`
	Size      MetaSize       `json:"size"`
	Scale     string         `json:"scale"`
	FrameTags []FrameTagItem `json:"frameTags"`
}

type SpriteAnimated struct {
	Frames   map[string]interface{} `json:"frames"`
	FrameArr []Frames
	Meta     Meta `json:"meta"`
}

func NewSpriteAnimated(fileName string) *SpriteAnimated {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(bufio.NewReader(jsonFile))
	var spriteAnimated SpriteAnimated
	json.Unmarshal(data, &spriteAnimated)

	var frames []Frames
	for key, value := range spriteAnimated.Frames {
		var frame Frames
		mapstructure.Decode(value, &frame)
		sliceKey := strings.Fields(key, "/n")
		fmt.Println(key, " --> ", sliceKey)
		frame.Name = sliceKey[0]
		frame.Action = sliceKey[1]
		// frame.Order = sliceKey[2]
		// fmt.Println("--> ", frame)
		frames = append(frames, frame)
	}
	spriteAnimated.FrameArr = frames
	// fmt.Println(spriteAnimated.FrameArr)

	return &spriteAnimated
}
