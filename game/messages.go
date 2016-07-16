package game

import tl "github.com/ilackarms/termloop"

type MessageType byte

const (
	Login = MessageType(0)
	Input = MessageType(1)
	Command = MessageType(2)
	LevelUpdate = MessageType(3)
)

type LoginMessage struct {
	Name string `json:"Name"`
}

type InputMessage struct {
	Event tl.Event `json:"Event"`
}

type CommandMessage struct {
	Text string `json:"Command"`
}

type LevelChangeMessage struct {
	LevelData levelData `json:"LevelData"`
}

type levelData struct {
	UUID string `json:"UUID"`
	Drawables map[string]drawableData `json:"Drawables"`
	Bg       tl.Cell `json:"Bg"`
	Offsetx  int `json:"Offsetx"`
	Offsety  int `json:"Offsety"`
}

type drawableData struct {
	Type tl.DrawableType `json:"Type"`
	Data []byte `json:"Data"`
}
