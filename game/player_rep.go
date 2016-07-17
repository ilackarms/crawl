package game

import (
	tl "github.com/ilackarms/termloop"
	"log"
	"github.com/emc-advanced-dev/pkg/errors"
	"encoding/json"
)

//player rep is the server's representation of the player.
//its position is meant to be updated through input commands sent to the server
//it should not be drawn and should not do anything on a tick
type PlayerRep struct {
	Name   string `json:"Name"`
	Entity *tl.Entity `json:"Entity"`
	PrevX  int `json:"PrevX"`
	PrevY  int `json:"PrevY"`
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

func (player *PlayerRep) ProcessEvent(event tl.Event) {
	if event.Type == tl.EventKey {
		// Is it a keyboard event?
		x, y := player.Entity.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.Entity.SetPosition(x + 1, y)
		case tl.KeyArrowLeft:
			player.Entity.SetPosition(x - 1, y)
		case tl.KeyArrowUp:
			player.Entity.SetPosition(x, y - 1)
		case tl.KeyArrowDown:
			player.Entity.SetPosition(x, y + 1)
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
}

const DrawableType_PlayerRep = tl.DrawableType("DrawableType_PlayerRep")

func DeserializePlayerRep(data []byte) (*PlayerRep, error) {
	var playerRep PlayerRep
	if err := json.Unmarshal(data, &playerRep); err != nil {
		return nil, errors.New("unmarshalling "+string(data)+" to playerRep", err)
	}
	return &playerRep, nil
}