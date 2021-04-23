package model

import (
	"fmt"
	"time"
)

var startTime int

const (
	// 1s = 1000ms
	oneSecond = 1000
)

type Fps struct {
	LastTime      int
	CurFps        int
	CurFrameCount int
	AllFrameCount int
	Visable       bool
}

func GetTimeStamp() int {
	return int(time.Now().UnixNano() / 1000000)
}

func NewFps(visable bool) *Fps {
	fps := Fps{
		LastTime:      startTime,
		CurFps:        0,
		CurFrameCount: 0,
		AllFrameCount: 0,
		Visable:       visable,
	}

	startTime = GetTimeStamp()

	return &fps
}

func (f *Fps) Tick() {
	f.CurFrameCount++
	f.AllFrameCount++

	curTime := GetTimeStamp()
	if curTime-f.LastTime > oneSecond {
		f.CurFps = f.CurFrameCount
		f.CurFrameCount = 0
		f.LastTime = curTime
		//console log to screen
		if f.Visable {
			fmt.Println(f.PrintFps())
			// fmt.Println(f.GetAverageFps())
			// fmt.Println(time.Unix(0, int64(f.GetStartTime()*1000000)))
			// fmt.Println(f.GetPlayTime()
		}
	}
}

func (f *Fps) GetCurFps() int {
	//when curTime-f.LastTime < 1000 use average fps
	if f.CurFps == 0 {
		return int(f.GetAverageFps())
	}

	return f.CurFps
}

func (f *Fps) PrintFps() string {
	if f.Visable {
		return fmt.Sprintf("fps=%d", f.CurFps)
	}

	return ""
}

func (f *Fps) SetVisable(visable bool) {
	f.Visable = visable
}

func (f *Fps) GetAverageFps() float32 {
	return float32(oneSecond * f.AllFrameCount / (f.GetPlayTime()))
}

func (f *Fps) GetStartTime() int {
	return startTime
}

func (f *Fps) GetPlayTime() int {
	return f.LastTime - startTime
}
