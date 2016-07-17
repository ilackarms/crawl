package game

import (
	tl "github.com/ilackarms/termloop"
	"net"
	"github.com/ilackarms/crawl/protocol"
	"strings"
)

const commandPrefix = "command: "

type Player struct {
	Name string
	entity *tl.Entity
	prevX  int
	prevY  int
	level  *tl.BaseLevel
	text *tl.Text
	server net.Conn
}

func NewPlayer(name string, entity *tl.Entity, server net.Conn) *Player {
	entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	return &Player{
		Name: name,
		entity: entity,
		server: server,
		text: tl.NewText(0, 0, "", tl.ColorWhite, tl.ColorBlack),
	}
}

func (player *Player) SetLevel(level *tl.BaseLevel) {
	player.level = level
	level.AddEntity(player.text)
}

func (player *Player) SetPosition(x, y int) {
	player.entity.SetPosition(x, y)
}

func (player *Player) GetUUID() string {
	return player.entity.GetUUID()
}

// Here, Draw simply tells the Entity ent to handle its own drawing.
// We don't need to do anything.
func (player *Player) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	x, y := player.entity.Position()
	player.level.SetOffset(screenWidth / 2 - x, screenHeight / 2 - y)
	player.text.SetPosition(x - len(player.text.Text())/2, y - 1 + screenHeight/2)
	player.entity.Draw(screen)
}

func (player *Player) Tick(event tl.Event) {
	player.prevX, player.prevY = player.entity.Position()
	if event.Type == tl.EventKey {
		// Is it a keyboard event?
		//x, y := player.entity.Position()
		currentText := player.text.Text()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			//player.entity.SetPosition(x + 1, y)
			fallthrough
		case tl.KeyArrowLeft:
			//player.entity.SetPosition(x - 1, y)
			fallthrough
		case tl.KeyArrowUp:
			//player.entity.SetPosition(x, y - 1)
			fallthrough
		case tl.KeyArrowDown:
			//player.entity.SetPosition(x, y + 1)
			player.sendEvent(event)
		case tl.KeyEnter:
			player.sendCommand(strings.TrimPrefix(currentText, commandPrefix))
			player.text.SetText(commandPrefix)
		case tl.KeySpace:
			player.text.SetText(currentText+" ")
		case tl.KeyBackspace:
			fallthrough
		case tl.KeyBackspace2:
			if len(currentText) <= len(commandPrefix) {
				player.text.SetText(commandPrefix)
			} else {
				player.text.SetText(currentText[:len(currentText)-1])
			}
		default:
			player.text.SetText(player.text.Text()+string(event.Ch))
		}
	}
}

func (player *Player) sendEvent(event tl.Event) error {
	message := InputMessage{Event: event}
	return protocol.SendMessage(player.server, message, Input)
}

func (player *Player) sendCommand(command string) error {
	message := CommandMessage{Text: Command}
	return protocol.SendMessage(player.server, message, Input)
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