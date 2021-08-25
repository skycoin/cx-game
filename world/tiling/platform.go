package tiling

import (
	"log"
)

const (
	platformTilingWidth int = 3
	platformTilingHeight int = 3
)

type PlatformTiling struct {}

func (t PlatformTiling) Count() int {
	return platformTilingWidth * platformTilingHeight
}

func (t PlatformTiling) Index(n DetailedNeighbours) int {
	if n.Left == Solid || n.Right == Solid { // 4-8
		if n.Left == Solid && n.Right == Self { return 4 }
		if n.Right == Solid && n.Left == Self { return 5 }
		if n.Left == Solid && n.Right == None { return 6 }
		if n.Right == Solid && n.Left == None { return 7 }
		if n.Left == Solid && n.Right == Solid { return 8 }
	} else { // 0-3
		if n.Right == Self && n.Left == None { return 0 }
		if n.Right == Self && n.Left == Self { return 1 }
		if n.Right == None && n.Left == Self { return 2 }
		if n.Right == None && n.Left == None { return 3 }
	}
	log.Fatalf("cannot find index for platform tiling\n%+v",n)
	return -1
}
