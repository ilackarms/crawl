package game

import (
	tl "github.com/ilackarms/termloop"
	"encoding/json"
)

type Level struct {
	*tl.BaseLevel
	PrevLevel []byte
	BeforeTick func(level *Level, ev tl.Event) `json:"-"`
	AfterTick func(level *Level, ev tl.Event) `json:"-"`
}

func (l *Level) Tick(ev tl.Event) {
	if l.BeforeTick != nil {
		l.BeforeTick(l, ev)
	}
	l.BaseLevel.Tick(ev)
	if l.AfterTick != nil {
		l.AfterTick(l, ev)
	}
}

func (l *Level) CacheLevel() {
	l.PrevLevel, _ = json.Marshal(l.Entities)
}

//returns true if l differs from PrevLevel
func (l *Level) Diff() bool {
	if len(l.PrevLevel) == 0 {
		return true
	}
	data, _ := json.Marshal(l.Entities)
	return string(data) != string(l.PrevLevel)
}