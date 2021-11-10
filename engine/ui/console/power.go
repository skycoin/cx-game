package console

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
)

func onOffToBool(onOff string) (bool,error) {
	if onOff == "on" {
		return true,nil
	} else if onOff == "off" {
		return false, nil
	}
	return false, errors.New("onOff is neither on nor off")
}

func Power(line string, ctx CommandContext) string {
	words := strings.Split(line," ")
	if len(words) < 4 {
		return "need 4 args"
	}
	x,err := strconv.Atoi(words[1])
	if err!=nil { return "parse error [x]" }
	y,err := strconv.Atoi(words[2])
	if err!=nil { return "parse error [y]" }
	onOff := words[3]
	isOn ,err := onOffToBool(onOff)
	if err!=nil { return err.Error() }

	ctx.World.Planet.TogglePower(x,y,isOn)
	return fmt.Sprintf("powered %d,%d %s", x,y, onOff)
}
