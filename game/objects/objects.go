package objects

import (
	"strings"
	tl "github.com/ilackarms/termloop"
)

const (
//	dungeonEntrance = `
// _____
//|*****|
//|*****|
//|*****|
//|*| |*|
//`
	dungeonEntrance = `
	0	1	2	3	4	5	6	7	8	9	A	B	C	D	E	F
U+246x	①	②	③	④	⑤	⑥	⑦	⑧	⑨	⑩	⑪	⑫	⑬	⑭	⑮	⑯
U+247x	⑰	⑱	⑲	⑳	⑴	⑵	⑶	⑷	⑸	⑹	⑺	⑻	⑼	⑽	⑾	⑿
U+248x	⒀	⒁	⒂	⒃	⒄	⒅	⒆	⒇	⒈	⒉	⒊	⒋	⒌	⒍	⒎	⒏
U+249x	⒐	⒑	⒒	⒓	⒔	⒕	⒖	⒗	⒘	⒙	⒚	⒛	⒜	⒝	⒞	⒟
U+24Ax	⒠	⒡	⒢	⒣	⒤	⒥	⒦	⒧	⒨	⒩	⒪	⒫	⒬	⒭	⒮	⒯
U+24Bx	⒰	⒱	⒲	⒳	⒴	⒵	Ⓐ	Ⓑ	Ⓒ	Ⓓ	Ⓔ	Ⓕ	Ⓖ	Ⓗ	Ⓘ	Ⓙ
U+24Cx	Ⓚ	Ⓛ	Ⓜ	Ⓝ	Ⓞ	Ⓟ	Ⓠ	Ⓡ	Ⓢ	Ⓣ	Ⓤ	Ⓥ	Ⓦ	Ⓧ	Ⓨ	Ⓩ
U+24Dx	ⓐ	ⓑ	ⓒ	ⓓ	ⓔ	ⓕ	ⓖ	ⓗ	ⓘ	ⓙ	ⓚ	ⓛ	ⓜ	ⓝ	ⓞ	ⓟ
U+24Ex	ⓠ	ⓡ	ⓢ	ⓣ	ⓤ	ⓥ	ⓦ	ⓧ	ⓨ	ⓩ	⓪	⓫	⓬	⓭	⓮	⓯
U+24Fx	⓰	⓱	⓲	⓳	⓴	⓵	⓶	⓷	⓸	⓹	⓺	⓻	⓼	⓽	⓾	⓿
─	━	│	┃	┄	┅	┆	┇	┈	┉	┊	┋	┌	┍	┎	┏
U+251x	┐	┑	┒	┓	└	┕	┖	┗	┘	┙	┚	┛	├	┝	┞	┟
U+252x	┠	┡	┢	┣	┤	┥	┦	┧	┨	┩	┪	┫	┬	┭	┮	┯
U+253x	┰	┱	┲	┳	┴	┵	┶	┷	┸	┹	┺	┻	┼	┽	┾	┿
U+254x	╀	╁	╂	╃	╄	╅	╆	╇	╈	╉	╊	╋	╌	╍	╎	╏
U+255x	═	║	╒	╓	╔	╕	╖	╗	╘	╙	╚	╛	╜	╝	╞	╟
U+256x	╠	╡	╢	╣	╤	╥	╦	╧	╨	╩	╪	╫	╬	╭	╮	╯
U+257x	╰	╱	╲	╳	╴	╵	╶	╷	╸	╹	╺	╻	╼	╽	╾	╿
U+2580	▀	Upper half block
U+2581	▁	Lower one eighth block
U+2582	▂	Lower one quarter block
U+2583	▃	Lower three eighths block
U+2584	▄	Lower half block
U+2585	▅	Lower five eighths block
U+2586	▆	Lower three quarters block
U+2587	▇	Lower seven eighths block
U+2588	█	Full block
U+2589	▉	Left seven eighths block
U+258A	▊	Left three quarters block
U+258B	▋	Left five eighths block
U+258C	▌	Left half block
U+258D	▍	Left three eighths block
U+258E	▎	Left one quarter block
U+258F	▏	Left one eighth block
U+2590	▐	Right half block
U+2591	░	Light shade
U+2592	▒	Medium shade
U+2593	▓	Dark shade
U+2594	▔	Upper one eighth block
U+2595	▕	Right one eighth block
U+2596	▖	Quadrant lower left
U+2597	▗	Quadrant lower right
U+2598	▘	Quadrant upper left
U+2599	▙	Quadrant upper left and lower left and lower right
U+259A	▚	Quadrant upper left and lower right
U+259B	▛	Quadrant upper left and upper right and lower left
U+259C	▜	Quadrant upper left and upper right and lower right
U+259D	▝	Quadrant upper right
U+259E	▞	Quadrant upper right and lower left
U+259F	▟	Quadrant upper right and lower left and lower right
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