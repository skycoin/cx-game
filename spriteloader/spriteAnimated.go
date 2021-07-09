package spriteloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
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
		sliceKey := regexp.MustCompile("[^\\s]+")
		sliceArr := sliceKey.FindAllString(key, -1)
		frame.Name = sliceArr[0]
		frame.Action = strings.Split(strings.ReplaceAll(sliceArr[1], "#", ""), ".")[0]
		if len(sliceArr) > 2 {
			i, err := strconv.Atoi(strings.Split(sliceArr[2], ".")[0])
			if err != nil {
				log.Fatal(err)
			}
			frame.Order = i
		} else {
			frame.Order = 0
		}
		// fmt.Println("----------------> ", key)
		// fmt.Println("name: ", frame.Name)
		// fmt.Println("Action: ", frame.Action)
		// fmt.Println("Order: ", frame.Order)
		frames = append(frames, frame)
	}
	spriteAnimated.FrameArr = frames

	return &spriteAnimated
}

func Draw(action string) {

}
