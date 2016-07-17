package game

import tl "github.com/ilackarms/termloop"

type Level struct {
	*tl.BaseLevel
	BeforeTick func(level Level, ev tl.Event) `json:"-"`
	AfterTick func(level Level, ev tl.Event) `json:"-"`
}

func (l *Level) Tick(ev tl.Event) {
	if l.BeforeTick != nil {
		l.BeforeTick(*l, ev)
	}
	l.BaseLevel.Tick(ev)
	if l.AfterTick != nil {
		l.AfterTick(*l, ev)
	}
}