package game

import (
	"encoding/json"
	tl "github.com/ilackarms/termloop"
)

type Level struct {
	*tl.BaseLevel
	Descriptions map[string]string `json:"Descriptions"` //WARNING: here's how to do it: fmt.Sprintf("%d,%d", x, y)
	PrevLevel    []byte
	BeforeTick   func(level *Level, ev tl.Event) `json:"-"`
	AfterTick    func(level *Level, ev tl.Event) `json:"-"`
}

func NewLevel(baseLevel *tl.BaseLevel) *Level {
	return &Level{
		Descriptions: make(map[string]string),
		BaseLevel:    baseLevel,
	}
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
