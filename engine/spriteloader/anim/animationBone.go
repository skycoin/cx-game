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

	fmt.Println("Skeleton: ", sheet.Skeleton)
	fmt.Println("Bones: ", sheet.Bones)
	fmt.Println("Slots: ", len(sheet.Slots))
	fmt.Println("Skins: ", sheet.Skins[0].Attachments)
	m, _ := sheet.Animations.(map[string]interface{})
	for k, v := range m {
		fmt.Println(k, "=>", v)
	}

	imgPath := "./assets/player/Robot.png"
	gpuTex := spriteloader.LoadTextureFromFileToGPU(imgPath)
	fmt.Println("gpuTex: ", gpuTex)

	// Bone.DrawBone(10, 10, 20, 20)

}
