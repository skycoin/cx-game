package world

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	DAY_DURATION = 10 //

	sunrise_length float32 = 0.15
	day_length     float32 = 0.4
	sunset_length  float32 = 0.15
	night_length   float32 = 0.3

	SUNRISE TimeOfDay = iota
	DAY
	SUNSET
	NIGHT
)

type TimeOfDay int

func (t TimeOfDay) ToString() string {
	switch t {
	case SUNRISE:
		return "sunrise"
	case DAY:
		return "day"
	case SUNSET:
		return "sunset"
	case NIGHT:
		return "night"
	}
	return "unknown"
}

type DayCycleController struct {
	dtime      float32 //current time of the day
	daysPassed uint32  //how many days passed since the beginning
	timeOfDay  TimeOfDay
	brightness float32
	LTG        LightTextureGenerator
}

func NewDayCycleController() *DayCycleController {
	dcc := DayCycleController{}
	return &dcc
}

func (dcc *DayCycleController) Advance(dt float32) {
	dcc.dtime += dt

	if dcc.dtime > DAY_DURATION {
		dcc.dtime -= DAY_DURATION
		dcc.daysPassed += 1
		fmt.Printf("Day %v started\n", dcc.daysPassed)
	}
	if dcc.dtime > DAY_DURATION {
		log.Fatalln("Not expected")
	}

	ct := dcc.dtime / DAY_DURATION

	if ct < 0.0 || ct > 1.0 {
		log.Fatalln("Out of possible range")
	}

	if ct < sunrise_length {
		//logic
		dcc.timeOfDay = SUNRISE
		return
	}
	ct -= sunrise_length
	if ct < day_length {
		dcc.timeOfDay = DAY
		//logic
		return
	}
	ct -= day_length
	if ct < sunrise_length {
		dcc.timeOfDay = SUNRISE
		//sunrise logic
		return
	}
	ct -= sunrise_length
	if ct < night_length {
		dcc.timeOfDay = NIGHT
		//night logic
		return
	}

	log.Fatalln("Should not be reached")
}

func (dcc DayCycleController) CurrentDay() uint32 {
	return dcc.daysPassed + 1
}

func (dcc DayCycleController) TimeOfDay() TimeOfDay {
	return dcc.timeOfDay
}

func (dcc *DayCycleController) SetTimeOfDay(tod TimeOfDay) {
	if tod == dcc.timeOfDay {
		return
	}
	dcc.timeOfDay = tod

	// fmt.Println("TIME OF DAY: ", dcc.timeOfDay.ToString())

	dcc.dtime = 0
	if tod == SUNRISE {
		return
	}
	dcc.dtime += sunrise_length
	if tod == DAY {
		return
	}
	dcc.dtime += day_length
	if tod == SUNSET {
		return
	}
	//night
	dcc.dtime += sunrise_length
}
func (dcc *DayCycleController) GetLightTexture() uint32 {
	return dcc.LTG.GenerateLightTexture(dcc.timeOfDay)
}

type LightTextureGenerator struct {
	lightTex uint32
}

func (ltg *LightTextureGenerator) GenerateLightTexture() uint32 {
	//generate 16x16 texture
	var lightTex uint32
	gl.GenTextures(1, &lightTex)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, lightTex)

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage3D(gl.TEXTURE_2D_ARRAY, 0, gl.RGBA, 16, 16, 16*16, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	return lightTex
}
