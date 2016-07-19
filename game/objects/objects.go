package objects

import (
	"strings"
	tl "github.com/ilackarms/termloop"
)

const (
	dungeonEntrance = `
▉ ▉ ▉ ▉
▉▉▉▉▉▉▉
▉░▉░▉░▉
▉▉▉▉▉▉▉
▉▉▉ ▉▉▉
`
)

func graphChars(objString string, color tl.Attr) tl.Canvas {
	lines := strings.Split(strings.Trim(objString, "\n"), "\n")
	var height int
	for i := range lines {
		line := []rune(lines[i])
		if len(line) > height {
			height = len(line)
		}
	}
	width := len(lines)
	canvas := make(tl.Canvas, height)
	for i := range canvas {
		canvas[i] = make([]tl.Cell, width)
	}
	for i := range lines {
		line := []rune(lines[i])
		for j := 0; j < len(line); j++ {
			canvas[j][i] = tl.Cell{
				Ch: line[j],
				Fg: color,
			}
		}
	}
	return canvas
}

type DungeonEntrance struct {
	*tl.Entity
	TargetLevelUUID string
}

func NewDungeonEntrance(x, y int, color tl.Attr, targetLevelUUID string) *DungeonEntrance {
	de := &DungeonEntrance{
		Entity: tl.NewEntityFromCanvas(x, y, graphChars(dungeonEntrance, color)),
		TargetLevelUUID: targetLevelUUID,
	}
	return de
}

//location of door relative to position
func (de *DungeonEntrance) TriggerPositions() []Position {
	x, y := de.Position()
	return []Position{
		Position{X: 3+x, Y: 0+y+de.Height-1},
	}
}