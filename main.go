package main

import tl "github.com/JoelOtter/termloop"

var input *tl.Text

func main() {
	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})
	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))
	game.Screen().SetLevel(level)
	player := Player{
		entity: tl.NewEntity(1, 1, 1, 1),
		level: level,
	}
	// Set the character at position (0, 0) on the entity.
	player.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	level.AddEntity(&player)
	_, height := game.Screen().Size()
	inputString := "input: "
	input = tl.NewText(0, height, inputString, tl.ColorWhite, tl.ColorBlack)
	level.AddEntity(input)
	game.Start()
}

type Player struct {
	entity *tl.Entity
	prevX  int
	prevY  int
	level  *tl.BaseLevel
}

// Here, Draw simply tells the Entity ent to handle its own drawing.
// We don't need to do anything.
func (player *Player) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	x, y := player.entity.Position()
	player.level.SetOffset(screenWidth / 2 - x, screenHeight / 2 - y)
	player.entity.Draw(screen)
}

func (player *Player) Tick(event tl.Event) {
	player.prevX, player.prevY = player.entity.Position()
	if event.Type == tl.EventKey {
		// Is it a keyboard event?
		x, y := player.entity.Position()
		currentText := input.Text()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.entity.SetPosition(x + 1, y)
		case tl.KeyArrowLeft:
			player.entity.SetPosition(x - 1, y)
		case tl.KeyArrowUp:
			player.entity.SetPosition(x, y - 1)
		case tl.KeyArrowDown:
			player.entity.SetPosition(x, y + 1)
		case tl.KeyEnter:
			input.SetText("input: ")
		case tl.KeySpace:
			input.SetText(currentText+" ")
		case tl.KeyBackspace:
			fallthrough
		case tl.KeyBackspace2:
			if len(currentText) <= len("input: ") {
				input.SetText("input: ")
			} else {
				input.SetText(currentText[:len(currentText)-1])
			}
		default:
			input.SetText(input.Text()+string(event.Ch))
		}
	}
}

func (player *Player) Size() (int, int) {
	return player.entity.Size()
}

func (player *Player) Position() (int, int) {
	return player.entity.Position()
}

func (player *Player) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		player.entity.SetPosition(player.prevX, player.prevY)
	}
}