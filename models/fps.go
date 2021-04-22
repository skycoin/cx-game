package model

import (
	"fmt"
	"time"
)

var StartTime int

type Fps struct {
	LastTime      int
	CurFps        int
	CurFrameCount int
	AllFrameCount int
}

func NewFps() *Fps {
	fps := Fps{}

	StartTime = int(time.Now().UnixNano() / 1000000)
	fps.LastTime = StartTime
	fps.CurFps = 0
	fps.CurFrameCount = 0
	fps.AllFrameCount = 0

	return &fps
}

func (f *Fps) Tick() {
	f.CurFrameCount++
	f.AllFrameCount++

	curTime := int(time.Now().UnixNano() / 1000000)
	if curTime-f.LastTime > 1000 {
		f.CurFps = f.CurFrameCount
		f.CurFrameCount = 0
		f.LastTime = curTime
		//console log to screen
		fmt.Println(f.CurFps)
		// fmt.Println(f.GetAverageFps())
		// fmt.Println(time.Unix(0, int64(f.GetStartTime()*1000000)))
		// fmt.Println(f.GetPlayTime())
	}
}

func (f *Fps) GetCurFps() int {
	//when curTime-f.LastTime < 1000 use average fps
	if f.CurFps == 0 {
		return int(f.GetAverageFps())
	}

	return f.CurFps
}

func (f *Fps) GetAverageFps() float32 {
	return float32(1000 * f.AllFrameCount / (f.LastTime - StartTime))
}

func (f *Fps) GetStartTime() int {
	return StartTime
}

func (f *Fps) GetPlayTime() int {
	return f.LastTime - StartTime
}
