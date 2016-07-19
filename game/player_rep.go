package game

import (
	tl "github.com/ilackarms/termloop"
	"log"
	"github.com/emc-advanced-dev/pkg/errors"
	"encoding/json"
	"github.com/ilackarms/crawl/game/objects"
)

//player rep is the server's representation of the player.
//its position is meant to be updated through input commands sent to the server
//it should not be drawn and should not do anything on a tick
type PlayerRep struct {
	Name   string `json:"Name"`
	Entity *tl.Entity `json:"Entity"`
	PrevX  int `json:"PrevX"`
	PrevY  int `json:"PrevY"`
	currentCommand string `json:"-"`
}

func NewPlayerRep(name string, entity *tl.Entity) *PlayerRep {
	entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	return &PlayerRep{
		Name: name,
		Entity: entity,
	}
}

func (player *PlayerRep) SetUUID(uuid string) {
	player.Entity.UUID = uuid
}

func (player *PlayerRep) GetUUID() string {
	return player.Entity.GetUUID()
}

func (player *PlayerRep) ProcessCommand(command CommandMessage) {
	player.currentCommand = command.Text
}

func (player *PlayerRep) targetCommand(x, y int) {
	log.Printf("TODO: Processing command %v at target %v,%v", player.currentCommand, x, y)
	player.currentCommand = ""
}

func (player *PlayerRep) cancelCommand() {
	log.Printf("TODO: Canceling command %v", player.currentCommand)
	player.currentCommand = ""
}

func (player *PlayerRep) ProcessInput(input InputMessage) {
	event := input.Event
	if event.Type == tl.EventMouse {
		switch event.Key { // If so, switch on the pressed key.
		case tl.MouseRelease:
			player.targetCommand(event.MouseX, event.MouseY)
		}
		return
	}
	if event.Type == tl.EventKey {
		// Is it a keyboard event?
		x, y := player.Entity.Position()
		player.PrevX = x
		player.PrevY = y
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.Entity.SetPosition(x + 1, y)
		case tl.KeyArrowLeft:
			player.Entity.SetPosition(x - 1, y)
		case tl.KeyArrowUp:
			player.Entity.SetPosition(x, y - 1)
		case tl.KeyArrowDown:
			player.Entity.SetPosition(x, y + 1)
		case tl.KeyEsc:
			player.cancelCommand()
		default:
			log.Fatalf("ERROR: unknown event %v", event)
		}
		return
	}
	log.Fatalf("ERROR: unknown event %v", event)
}

func (player *PlayerRep) Draw(screen *tl.Screen) {
	player.Entity.Draw(screen)
}

func (player *PlayerRep) Tick(event tl.Event) {
	player.PrevX, player.PrevY = player.Entity.Position()
}

func (player *PlayerRep) Size() (int, int) {
	return player.Entity.Size()
}

func (player *PlayerRep) Position() (int, int) {
	return player.Entity.Position()
}

func (player *PlayerRep) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		player.Entity.SetPosition(player.PrevX, player.PrevY)
	}
	if trigger, ok := collision.(objects.Trigger); ok {
		x, y := player.Position()
		for _, triggerPosition := range trigger.TriggerPositions() {
			if x == triggerPosition.X && y == triggerPosition.Y {
				if dungeonEntrance, ok := trigger.(*objects.DungeonEntrance); ok {
					log.Printf("entering %v", dungeonEntrance)
				}
			}
		}
		player.Entity.SetPosition(player.PrevX, player.PrevY)
	}
}

const DrawableType_PlayerRep = tl.DrawableType("DrawableType_PlayerRep")

func DeserializePlayerRep(data []byte) (*PlayerRep, error) {
	var playerRep PlayerRep
	if err := json.Unmarshal(data, &playerRep); err != nil {
		return nil, errors.New("unmarshalling "+string(data)+" to playerRep", err)
	}
	return &playerRep, nil
}