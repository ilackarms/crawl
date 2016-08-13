package game

import tl "github.com/ilackarms/termloop"

type MessageType byte

func (t MessageType) GetByte() byte {
	return byte(t)
}

const (
	Login       = MessageType(1)
	Input       = MessageType(2)
	Command     = MessageType(3)
	LevelUpdate = MessageType(4)
)

type LoginMessage struct {
	Name string `json:"Name"`
	UUID string `json:"UUID"`
}

type InputMessage struct {
	Event       tl.Event `json:"Event"`
	CustomEvent func()   `json:"-"`
}

type CommandMessage struct {
	Text string `json:"Command"`
}

type LevelChangeMessage struct {
	LevelData levelData `json:"LevelData"`
}

type levelData struct {
	UUID      string         `json:"UUID"`
	Drawables []drawableData `json:"Drawables"`
	Bg        tl.Cell        `json:"Bg"`
	Offsetx   int            `json:"Offsetx"`
	Offsety   int            `json:"Offsety"`
}

type drawableData struct {
	Type tl.DrawableType `json:"Type"`
	Data []byte          `json:"Data"`
}
