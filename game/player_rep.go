package game

import (
	"encoding/json"
	"github.com/emc-advanced-dev/pkg/errors"
	tl "github.com/ilackarms/termloop"
	"log"
)

//player rep is the server's representation of the player.
//its position is meant to be updated through input commands sent to the server
//it should not be drawn and should not do anything on a tick
type PlayerRep struct {
	Name           string      `json:"Name"`
	Entity         *tl.Entity  `json:"Entity"`
	PrevX          int         `json:"PrevX"`
	PrevY          int         `json:"PrevY"`
	currentCommand string      `json:"-"`
	Iq             *inputQueue `json:"-"`
	W              *World      `json:"-"`
}

func NewPlayerRep(name string, entity *tl.Entity, w *World) *PlayerRep {
	entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '옷'})
	return &PlayerRep{
		Name:   name,
		Entity: entity,
		W:      w,
		Iq:     &inputQueue{},
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
	player.Iq.Push(input)
}

func (player *PlayerRep) Draw(screen *tl.Screen) {
	player.Entity.Draw(screen)
}

func (player *PlayerRep) Tick(event tl.Event) {
	player.PrevX, player.PrevY = player.Entity.Position()
	if player.Iq != nil && player.Iq.hasNext() { //iq is always nil on client side
		input := player.Iq.pop()
		if input.CustomEvent != nil {
			input.CustomEvent()
			return
		}
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
			switch event.Key { // If so, switch on the pressed key.
			case tl.KeyArrowRight:
				player.Entity.SetPosition(x+1, y)
			case tl.KeyArrowLeft:
				player.Entity.SetPosition(x-1, y)
			case tl.KeyArrowUp:
				player.Entity.SetPosition(x, y-1)
			case tl.KeyArrowDown:
				player.Entity.SetPosition(x, y+1)
			case tl.KeyEsc:
				player.cancelCommand()
			default:
				log.Fatalf("ERROR: unknown event %v", event)
			}
			x, y = player.Entity.Position()
			log.Printf("current position %v,%v", x, y)
			return
		}
		log.Fatalf("ERROR: unknown event %v", event)
	}
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
	} else if trigger, ok := collision.(*Trigger); ok {
		x, y := player.Position()
		if callback, ok := trigger.TriggerPositions()[Position{X: x, Y: y}]; ok {
			log.Printf("trigger: %v,%v", x, y)
			callback(player)
		} else {
			log.Printf("~trigger: %v,%v", x, y)
			player.Entity.SetPosition(player.PrevX, player.PrevY)
		}
	} else {
		//log.Printf("collisino: %v", collision)
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

type inputQueue struct {
	inputs []*InputMessage
}

func (iq *inputQueue) Push(input InputMessage) {
	iq.inputs = append(iq.inputs, &input)
}

func (iq *inputQueue) pop() InputMessage {
	input := iq.inputs[0]
	iq.inputs = iq.inputs[1:]
	return *input
}

func (iq *inputQueue) hasNext() bool {
	return len(iq.inputs) > 0
}
