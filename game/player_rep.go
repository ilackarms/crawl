package game

import (
	tl "github.com/ilackarms/termloop"
	"log"
)

//player rep is the server's representation of the player.
//its position is meant to be updated through input commands sent to the server
//it should not be drawn and should not do anything on a tick
type PlayerRep struct {
	Name string
	entity *tl.Entity
	prevX  int
	prevY  int
}

func NewPlayerRep(name string, entity *tl.Entity) *PlayerRep {
	return &PlayerRep{
		Name: name,
		entity: entity,
	}
}

func (player *PlayerRep) GetUUID() string {
	return player.entity.GetUUID()
}

func (player *PlayerRep) Draw(screen *tl.Screen) {
	//don't draw on server side
}

func (player *PlayerRep) Tick(event tl.Event) {
	player.prevX, player.prevY = player.entity.Position()
}

func (player *PlayerRep) Size() (int, int) {
	return player.entity.Size()
}

func (player *PlayerRep) Position() (int, int) {
	return player.entity.Position()
}

func (player *PlayerRep) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		player.entity.SetPosition(player.prevX, player.prevY)
	}
}