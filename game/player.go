package game

import (
	tl "github.com/ilackarms/termloop"
	"net"
	"github.com/ilackarms/crawl/protocol"
)

type Player struct {
	Name   string
	entity *tl.Entity
	prevX  int
	prevY  int
	level  *Level
	Text   *tl.Text
	server net.Conn
}

func NewPlayer(name string, entity *tl.Entity, server net.Conn) *Player {
	entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	return &Player{
		Name: name,
		entity: entity,
		server: server,
		Text: tl.NewText(0, 0, "", tl.ColorWhite, tl.ColorBlack),
	}
}

func (player *Player) SetLevel(level *Level) {
	player.level = level
	level.AddEntity(player.Text)
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
	player.Text.SetPosition(x - len(player.Text.GetText())/2, y - 1 + screenHeight/2)
	player.entity.Draw(screen)
}

func (player *Player) Tick(event tl.Event) {
	player.prevX, player.prevY = player.entity.Position()
	if event.Type == tl.EventMouse {
		switch event.Key {
		case tl.MouseRelease:
			player.sendEvent(event)
		}
	}
	if event.Type == tl.EventKey {
		// Is it a keyboard event?
		currentText := player.Text.GetText()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyEsc:
			fallthrough
		case tl.KeyArrowRight:
			fallthrough
		case tl.KeyArrowLeft:
			fallthrough
		case tl.KeyArrowUp:
			fallthrough
		case tl.KeyArrowDown:
			player.sendEvent(event)
		case tl.KeyEnter:
			player.sendCommand(currentText)
			player.Text.SetText("")
		case tl.KeySpace:
			player.Text.SetText(currentText+" ")
		case tl.KeyBackspace:
			fallthrough
		case tl.KeyBackspace2:
			if len(currentText) <= 1{
				player.Text.SetText("")
			} else {
				player.Text.SetText(currentText[:len(currentText)-1])
			}
		default:
			player.Text.SetText(player.Text.GetText()+string(event.Ch))
		}
	}
}

func (player *Player) sendEvent(event tl.Event) error {
	message := InputMessage{Event: event}
	return protocol.SendMessage(player.server, message, Input.GetByte())
}

func (player *Player) sendCommand(command string) error {
	message := CommandMessage{Text: command}
	return protocol.SendMessage(player.server, message, Command.GetByte())
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