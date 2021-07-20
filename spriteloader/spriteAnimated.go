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
	"runtime"

	"github.com/mitchellh/mapstructure"
)

func init() {
	runtime.LockOSThread()
}

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

func NewSpriteAnimated(fileName string) *SpriteAnimated {
	spriteAnimated := SpriteAnimated{}
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
	// assuming uniform dimensions
	cols := spriteAnimated.Meta.Size.W / spriteAnimated.FrameArr[0].Frame.W
	rows := spriteAnimated.Meta.Size.H / spriteAnimated.FrameArr[0].Frame.H
	spriteAnimated.spriteSheetId =
		LoadSpriteSheetByColRow("./assets/"+spriteAnimated.Meta.Image,rows,cols)
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


type Action struct {
	SpritesheetID SpritesheetID
	Frames []ActionFrame
	Time float32
	FrameIndex int
}
type ActionFrame struct {
	SpriteID SpriteID
	Duration float32
}

func (spriteAnimated SpriteAnimated) Action(name string) Action {
	framess := filterByAction(name, spriteAnimated.FrameArr)
	actionframes := make([]ActionFrame,len(framess))
	for idx,frames := range framess {
		actionframes[idx] = ActionFrame {
			SpriteID: LoadSprite(
				spriteAnimated.spriteSheetId,frames.Name,
				frames.Frame.X / frames.Frame.W,
				frames.Frame.Y / frames.Frame.H ),
			Duration: float32(frames.Duration)/1000,
		}
	}
	return Action {
		SpritesheetID: spriteAnimated.spriteSheetId,
		Frames: actionframes,
	}
}

func (action *Action) SpriteID() SpriteID {
	return action.Frame().SpriteID
}

func (action *Action) Frame() ActionFrame {
	return action.Frames[action.FrameIndex]
}

func (action *Action) Update(dt float32) {
	// accumulate time
	action.Time += dt
	// consume time with frames
	for action.Time > action.Frame().Duration {
		action.Time -= action.Frame().Duration
		action.FrameIndex = (action.FrameIndex+1)%len(action.Frames)
	}
}
