package game

import (
	"encoding/json"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/ilackarms/crawl/protocol"
	tl "github.com/ilackarms/termloop"
	"log"
	"net"
)

type Client struct {
	PlayerRep *PlayerRep `json:"PlayerRep"`
	conn      net.Conn   `json:"-"`
}

type World struct {
	CurrentLevel string
	Levels       map[string]tl.Level
	Clients      map[string]*Client
	G            *tl.Game
}

func NewWorld() *World {
	return &World{
		CurrentLevel: "",
		Levels:       make(map[string]tl.Level),
		Clients:      make(map[string]*Client),
		G:            tl.NewGame(),
	}
}

func (w *World) StartGame() {
	go w.G.StartServerMode()
}

func (w *World) AddLevel(level *Level) {
	level.AfterTick = w.syncLevel
	w.Levels[level.UUID] = level
}

func (w *World) SetLevel(levelUUID string) {
	w.CurrentLevel = levelUUID
	level, ok := w.Levels[w.CurrentLevel]
	if !ok {
		log.Fatalf("invalid level selected: %v, available levels: %v", levelUUID, w.Levels)
	}
	w.G.Screen().SetLevel(level)
}

func (w *World) HandleClient(conn net.Conn) {
	var clientUUID string
	for {
		message, messageType, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Fatalf("ERROR: reading client message: %v", err)
		}
		log.Printf("recieved: %v %v %v", string(message), messageType, err)
		switch messageType {
		case Login.GetByte():
			var login LoginMessage
			if err := json.Unmarshal(message, &login); err != nil {
				log.Fatalf("ERROR: client login: %v", err)
			}
			playerRep := NewPlayerRep(login.Name, tl.NewEntity(0, 0, 1, 1), w)
			client := &Client{
				PlayerRep: playerRep,
				conn:      conn,
			}
			clientUUID = login.UUID
			client.PlayerRep.SetUUID(login.UUID)
			w.Levels[w.CurrentLevel].AddEntity(playerRep)
			w.Clients[clientUUID] = client
		case Input.GetByte():
			if clientUUID == "" {
				log.Fatalf("ERROR: client has not logged in yet")
			}
			var input InputMessage
			if err := json.Unmarshal(message, &input); err != nil {
				log.Fatalf("ERROR: client input: %v", err)
			}
			w.Clients[clientUUID].PlayerRep.ProcessInput(input)
		case Command.GetByte():
			if clientUUID == "" {
				log.Fatalf("ERROR: client has not logged in yet")
			}
			var command CommandMessage
			if err := json.Unmarshal(message, &command); err != nil {
				log.Fatalf("ERROR: client command: %v", err)
			}
			w.Clients[clientUUID].PlayerRep.ProcessCommand(command)
		}
	}
}

func (w *World) syncLevel(level *Level, ev tl.Event) {
	if level.Diff() {
		for _, client := range w.Clients {
			if err := func() error {
				levelData, err := SerializeLevel(*level)
				if err != nil {
					return errors.New("serializing level", err)
				}
				if err := protocol.SendMessage(client.conn, LevelChangeMessage{LevelData: levelData}, LevelUpdate.GetByte()); err != nil {
					return errors.New("writing level data to client", err)
				}
				return nil
			}(); err != nil {
				log.Fatalf("ERROR: failed syncing level with client: %v", err)
			}
		}
		level.CacheLevel()
	}
}
