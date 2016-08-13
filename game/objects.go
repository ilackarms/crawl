package game

import (
	tl "github.com/ilackarms/termloop"
	"strings"
)

const (
	DungeonEntrance = `
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
