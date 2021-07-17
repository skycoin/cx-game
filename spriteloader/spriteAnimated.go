package spriteloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/mitchellh/mapstructure"
	"github.com/skycoin/cx-game/render"
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
	Frames        map[string]interface{} `json:"frames"`
	FrameArr      []Frames
	Meta          Meta `json:"meta"`
	spriteSheetId SpritesheetID
}

var spriteAnimated SpriteAnimated
var spriteId SpriteID
var stopPlay chan bool

func NewSpriteAnimated(fileName string, lwin *render.Window) *SpriteAnimated {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(bufio.NewReader(jsonFile))
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
		frames = append(frames, frame)
	}
	spriteAnimated.FrameArr = frames
	// load sprite
	InitSpriteloader(lwin)
	// spriteAnimated.spriteSheetId = LoadSpriteSheetByFrames("./assets/"+spriteAnimated.Meta.Image, spriteAnimated.FrameArr)
	spriteAnimated.spriteSheetId = LoadSpriteSheetByColRow("./assets/"+spriteAnimated.Meta.Image, 5, 7)
	// sorting frame by Action and Order
	sort.SliceStable(spriteAnimated.FrameArr, func(i, j int) bool {
		frI, frJ := spriteAnimated.FrameArr[i], spriteAnimated.FrameArr[j]
		switch {
		case frI.Action != frJ.Action:
			return frI.Action < frJ.Action
		default:
			return frI.Order < frJ.Order
		}
	})
	// fmt.Println(spriteAnimated.FrameArr)
	return &spriteAnimated
}

func filterByAction(action string, frames []Frames) []Frames {
	result := []Frames{}
	for i := range frames {
		if frames[i].Action == action {
			result = append(result, frames[i])
		}
	}
	return result
}

func (spriteAnimated *SpriteAnimated) Play(glwindow *glfw.Window, action string) {
	frames := filterByAction(action, spriteAnimated.FrameArr)
	stopPlay = make(chan bool)
	j := 0
	for {
		select {
		default:
			time.Sleep(time.Duration(frames[j].Duration) * time.Millisecond)
			LoadSprite(spriteAnimated.spriteSheetId, frames[j].Name, frames[j].Frame.X, frames[j].Frame.Y)
			spriteId := GetSpriteIdByName(frames[j].Name)
			fmt.Println("frames. ", frames[j].Action, " - ", frames[j].Order, " X.", frames[j].Frame.X, " Y.", frames[j].Frame.Y, " H.", frames[j].Frame.H, " W.", frames[j].Frame.W)
			if err := gl.Init(); err != nil {
				panic(err)
			}
			gl.ClearColor(1, 1, 1, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			DrawSpriteQuad(0, 0, 3, 2, spriteId)
			glwindow.SwapBuffers()
			glfw.PollEvents()
			j++
			if j == len(frames) {
				j = 0
			}
		case <-stopPlay:
			close(stopPlay)
			return
		}
	}
}
