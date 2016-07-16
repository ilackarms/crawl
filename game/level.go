package game

import tl "github.com/ilackarms/termloop"

type Level struct {
	tl.BaseLevel
	Callback func(level Level)
}

func (l *Level) Tick(ev tl.Event) {
	l.BaseLevel.Tick(ev)
	l.Callback(l)
}