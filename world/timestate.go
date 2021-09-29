package world

import (
	"log"

	"github.com/skycoin/cx-game/cxmath/math32"
)

type TimeState struct {
	TickCount      int
	TicksPerSecond int
	//to offset cycle timej
	//in seconds
	SecondsPassed int
	DayDuration   int
	TimeOffset    float32
	SunriseLength float32
	DayLength     float32
	SunsetLength  float32
	NightLength   float32
	LTG           LightTextureGenerator
}

func NewTimeState() TimeState {
	newTimeState := TimeState{
		TickCount:      0,
		SecondsPassed:  0,
		DayDuration:    10,
		TicksPerSecond: 30,
		SunriseLength:  0.15,
		DayLength:      0.7,
		SunsetLength:   0.15,
		NightLength:    0.0,

		LTG: NewLightTextureGenerator(),
	}

	//start at noon
	newTimeState.TimeOffset = newTimeState.SunriseLength
	if !checkTimeState(newTimeState) {
		log.Fatalln("Wrong timestate")
	}

	return newTimeState

}

func checkTimeState(timeState TimeState) bool {
	sum := timeState.DayLength + timeState.SunriseLength + timeState.SunsetLength + timeState.NightLength

	if sum != 1 {
		return false
	}
	return true
}

func (t *TimeState) Advance() {

	if t.TickCount%t.TicksPerSecond == 0 {
		t.SecondsPassed += 1

		ttime := t.SecondsPassed % t.DayDuration
		//where are we in a day between 0 and 1
		dtime := float32(ttime) / float32(t.DayDuration)

		//start from day
		lightV := t.calcLightValue(dtime)
		t.LTG.GenerateLightTexture(lightV)

	}
	t.TickCount++
}

func (t *TimeState) DaysPassed() int {
	return int(t.SecondsPassed / t.DayDuration)
}
func (t *TimeState) calcLightValue(ttime float32) float32 {
	if !(ttime >= 0 && ttime <= 1) {
		log.Fatalf("Should be between 0 and 1, have: %v\n", ttime)
	}

	ttime = math32.Mod(ttime+t.TimeOffset, 1)

	if ttime < t.SunriseLength {
		return ttime / t.SunriseLength
	}
	ttime -= t.SunriseLength
	if ttime < t.DayLength {
		return 1.0
	}
	ttime -= t.DayLength

	if ttime < t.SunsetLength {
		return 1.0 - ttime/t.SunsetLength
	}
	ttime -= t.SunsetLength
	if ttime <= t.NightLength {
		return 0.0
	}

	log.Fatalf("Should not reach this code, lightValue: %v\n", ttime)
	return 1.0
}
