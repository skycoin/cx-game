package anim

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/skycoin/cx-game/engine/spriteloader"
	animjson "github.com/skycoin/cx-game/engine/spriteloader/anim/json"
)

func LoadAnimationBoneFromJSON(fname string) {
	buff, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("could not find animation spritesheet at %s", fname)
	}
	sheet := animjson.UnmarshalAnimatedBoneSpriteSheet(buff)

	fmt.Println("sheet: ", sheet.Skeleton.Hash)
	imgPath := "./assets/player/Robot.png"
	gpuTex := spriteloader.LoadTextureFromFileToGPU(imgPath)
	fmt.Println("gpuTex: ", gpuTex)

}
