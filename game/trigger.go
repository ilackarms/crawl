package game

import (
	tl "github.com/ilackarms/termloop"
)

type Trigger struct {
	*tl.Entity
	Triggers map[Position]func(player *PlayerRep) `json:"-"`
}

func NewTrigger(x, y int, triggers map[Position]func(player *PlayerRep), objString string, color tl.Attr) *Trigger {
	de := &Trigger{
		Entity: tl.NewEntityFromCanvas(x, y, graphChars(objString, color)),
		Triggers: triggers,
	}
	return de
}

//location of door relative to position
func (t *Trigger) TriggerPositions() map[Position]func(player *PlayerRep) {
	x, y := t.Position()
	triggers := make(map[Position]func(player *PlayerRep))
	for position, callback := range t.Triggers {
		triggers[Position{X: position.X+x, Y: position.Y+y+t.Height-1}] = callback
	}
	return triggers
}