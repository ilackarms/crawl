package game

import (
	"time"
	"math/rand"
	tl "github.com/ilackarms/termloop"
)

type Point struct {
	x int
	y int
	p *Point
}

func (p *Point) Opposite() *Point {
	if p.x != p.p.x {
		return &Point{x: p.x + (p.x - p.p.x), y: p.y, p: p}
	}
	if p.y != p.p.y {
		return &Point{x: p.x, y: p.y + (p.y - p.p.y), p: p}
	}
	return nil
}


func adjacents(point *Point, maze [][]rune) []*Point {
	res := make([]*Point, 0)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if (i == 0 && j == 0) || (i != 0 && j != 0) {
				continue
			}
			if !isInMaze(point.x+i, point.y+j, len(maze), len(maze[0])) {
				continue
			}
			if maze[point.x+i][point.y+j] == '*' {
				res = append(res, &Point{point.x + i, point.y + j, point})
			}
		}
	}
	return res
}

func isInMaze(x, y int, w, h int) bool {
	return x >= 0 && x < w &&
	y >= 0 && y < h
}

// Generates a maze using Prim's Algorithm
// https://en.wikipedia.org/wiki/Maze_generation_algorithm#Randomized_Prim.27s_algorithm
func generateMaze(w, h int) [][]rune {
	maze := make([][]rune, w)
	for row := range maze {
		maze[row] = make([]rune, h)
		for ch := range maze[row] {
			maze[row][ch] = '*'
		}
	}
	rand.Seed(time.Now().UnixNano())
	point := &Point{x: rand.Intn(w), y: rand.Intn(h)}
	maze[point.x][point.y] = 'S'
	var last *Point
	walls := adjacents(point, maze)
	for len(walls) > 0 {
		rand.Seed(time.Now().UnixNano())
		wall := walls[rand.Intn(len(walls))]
		for i, w := range walls {
			if w.x == wall.x && w.y == wall.y {
				walls = append(walls[:i], walls[i+1:]...)
				break
			}
		}
		opp := wall.Opposite()
		if isInMaze(opp.x, opp.y, w, h) && maze[opp.x][opp.y] == '*' {
			maze[wall.x][wall.y] = '.'
			maze[opp.x][opp.y] = '.'
			walls = append(walls, adjacents(opp, maze)...)
			last = opp
		}
	}
	maze[last.x][last.y] = 'L'
	bordered := make([][]rune, len(maze)+2)
	for r := range bordered {
		bordered[r] = make([]rune, len(maze[0])+2)
		for c := range bordered[r] {
			if r == 0 || r == len(maze)+1 || c == 0 || c == len(maze[0])+1 {
				bordered[r][c] = '*'
			} else {
				bordered[r][c] = maze[r-1][c-1]
			}
		}
	}
	return bordered
}

func NewDungeonLevel(w, h int) *Level {
	bg := tl.RgbTo256Color(25, 25, 25)
	baseLevel := tl.NewBaseLevel(tl.Cell{
		Bg: bg,
		Fg: tl.ColorBlack,
		Ch: ' ',
	})
	maze := generateMaze(w, h)
	for i, row := range maze {
		for j, el := range row {
			if el == '*' && !(i == w/2 && j == h/2) {
				wall := tl.NewRectangle(i - w/2, j - h/2, 1, 1, bg)
				wall.Ch = '▓'
				wall.Fg = tl.RgbTo256Color(100, 110, 100)
				baseLevel.AddEntity(wall)
			}
		}
	}

	return &Level{
		BaseLevel: baseLevel,
	}
}