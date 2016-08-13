package game

import (
	"github.com/ilackarms/crawl/protocol"
	tl "github.com/ilackarms/termloop"
	"net"
)

type Player struct {
	Name                 string
	entity               *tl.Entity
	prevX                int
	prevY                int
	level                *Level
	InputText            *tl.Text
	descriptionTextLine1 *tl.Text
	descriptionTextLine2 *tl.Text
	descriptionTextLine3 *tl.Text
	description          string
	server               net.Conn
}

func NewPlayer(name string, entity *tl.Entity, server net.Conn) *Player {
	entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	return &Player{
		Name:                 name,
		entity:               entity,
		server:               server,
		InputText:            tl.NewText(0, 0, "", tl.ColorWhite, tl.ColorBlack),
		descriptionTextLine1: tl.NewText(0, 0, "", tl.ColorWhite, tl.RgbTo256Color(0, 0, 0)),
		descriptionTextLine2: tl.NewText(0, 0, "", tl.ColorWhite, tl.RgbTo256Color(0, 0, 0)),
		descriptionTextLine3: tl.NewText(0, 0, "", tl.ColorWhite, tl.RgbTo256Color(0, 0, 0)),
	}
}

func (player *Player) SetDescription(text string) {
	player.description = text
}

func (player *Player) SetLevel(level *Level) {
	player.level = level
	level.AddEntity(player.InputText)
	level.AddEntity(player.descriptionTextLine1)
	level.AddEntity(player.descriptionTextLine2)
	level.AddEntity(player.descriptionTextLine3)
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
	player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	player.InputText.SetPosition(x-len(player.InputText.GetText())/2, y+screenHeight/2)
	player.drawDescription(x, y, screenWidth, screenHeight)
	player.entity.Draw(screen)
}

func (player *Player) drawDescription(x, y, screenWidth, screenHeight int) {
	var d1, d2, d3 string
	d1 = player.description
	if len(player.description) > screenWidth {
		d1 = player.description[:screenWidth]
		d2 = player.description[screenWidth:]
	}
	if len(player.description) > screenWidth*2 {
		d2 = player.description[screenWidth : screenWidth*2]
		d3 = player.description[screenWidth*2:]
	}
	player.descriptionTextLine1.SetText(d1)
	player.descriptionTextLine2.SetText(d2)
	player.descriptionTextLine3.SetText(d3)
	player.descriptionTextLine1.SetPosition(x-screenWidth/2, y-3+screenHeight/2)
	player.descriptionTextLine2.SetPosition(x-screenWidth/2, y-2+screenHeight/2)
	player.descriptionTextLine3.SetPosition(x-screenWidth/2, y-1+screenHeight/2)
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
		currentText := player.InputText.GetText()
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
			player.InputText.SetText("")
		case tl.KeySpace:
			player.InputText.SetText(currentText + " ")
		case tl.KeyBackspace:
			fallthrough
		case tl.KeyBackspace2:
			if len(currentText) <= 1 {
				player.InputText.SetText("")
			} else {
				player.InputText.SetText(currentText[:len(currentText)-1])
			}
		default:
			player.InputText.SetText(player.InputText.GetText() + string(event.Ch))
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
