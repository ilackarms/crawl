package objects

import (
	"strings"
	tl "github.com/ilackarms/termloop"
)

const (
	dungeonEntrance = `
 _____
|*****|
|*****|
|*****|
|*| |*|
`
)

func graphChars(objString string, color tl.Attr) tl.Canvas {
	lines := strings.Split(strings.Trim(objString, "\n"), "\n")
	var height int
	for _, line := range lines {
		if len(line) > height {
			height = len(line)
		}
	}
	width := len(lines)
	canvas := make(tl.Canvas, height)
	for i := range canvas {
		canvas[i] = make([]tl.Cell, width)
	}
	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			canvas[j][i] = tl.Cell{
				Ch: rune(line[j]),
				Fg: color,
			}
		}
	}
	return canvas
}

type DungeonEntrance struct {
	*tl.Entity
}

func NewDungeonEntrance(x, y int, color tl.Attr) *DungeonEntrance {
	de := &DungeonEntrance{
		Entity: tl.NewEntityFromCanvas(x, y, graphChars(dungeonEntrance, color)),
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