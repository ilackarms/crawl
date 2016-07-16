package server

import (
	"net"
	tl "github.com/ilackarms/termloop"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"github.com/ilackarms/crawl/game"
	"github.com/ilackarms/crawl/protocol"
	"encoding/json"
)

type Client struct {
	PlayerRep     game.PlayerRep `json:"PlayerRep"`
	conn       net.Conn `json:"-"`
}

var (
	clients map[string]*Client
	levels map[string]tl.Level
	currentLevel string
	g *tl.Game
)

func init() {
	clients = make(map[string]*Client)
	levels = make(map[string]tl.Level)
}

func Start() {
	//start game
	g = startGame()

	//start multiplayer server
	if err := startServer(); err != nil {
		log.Fatal(errors.New("starting server", err))
	}

	syncLevel := func(level game.Level, ev tl.Event) {
		for _, client := range clients {
			if err := func() error {
				levelData, err := game.SerializeLevel(level)
				if err != nil {
					return errors.New("serializing level", err)
				}
				if err := protocol.SendMessage(client.conn, levelData, game.LevelUpdate); err != nil {
					return errors.New("writing level data to client", err)
				}
				return nil
			}(); err != nil {
				log.Printf("ERROR: failed syncing level with client: %v", err)
			}
		}
	}

	//test - create  & set the current level
	level1 := game.Level{
		BaseLevel: tl.NewBaseLevel(tl.Cell{
			Bg: tl.ColorGreen,
			Fg: tl.ColorBlack,
			Ch: 'v',
		}),
		AfterTick: syncLevel,
	}
	level1.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))
	levels[level1.UUID] = level1
	currentLevel = level1.UUID
	g.Screen().SetLevel(levels[currentLevel])
	select {}
	//test
}

func startGame() *tl.Game {
	g := tl.NewGame()
	go g.StartServerMode()
	return g
}

func startServer() error {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		return errors.New("starting listener", err)
	}
	log.Printf("listening on :9000")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("ERROR: failed to accept connection: %v", err)
				continue
			}
			go handle(conn)
		}
	}()
	return nil
}

func handle(conn net.Conn) {
	for {
		var clientUUID string
		message, messageType, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Printf("ERROR: reading client message: %v", err)
			continue
		}
		log.Printf("recieved: %v %v %v", message, messageType, err)
		switch messageType {
		case game.Login:
			var login game.LoginMessage
			if err := json.Unmarshal(message, &login); err != nil {
				log.Printf("ERROR: client login: %v", err)
				continue
			}
			playerRep := game.NewPlayerRep(login.Name, tl.NewEntity(1, 1, 1, 1))
			client := Client{
				PlayerRep: playerRep,
				conn: conn,
			}
			clientUUID = client.PlayerRep.GetUUID()
			levels[currentLevel].AddEntity(playerRep)
			clients[clientUUID] = client
		case game.Input:
			if clientUUID == "" {
				log.Printf("ERROR: client has not logged in yet")
				continue
			}
			var input game.InputMessage
			if err := json.Unmarshal(message, &input); err != nil {
				log.Printf("ERROR: client input: %v", err)
				continue
			}
			clients[clientUUID].PlayerRep.ProcessEvent(input.Event)
		case game.Command:
			if clientUUID == "" {
				log.Printf("ERROR: client has not logged in yet")
				continue
			}
			var command game.CommandMessage
			if err := json.Unmarshal(message, &command); err != nil {
				log.Printf("ERROR: client command: %v", err)
				continue
			}
			log.Printf("TODO: handle this command: %v", command.Text)
		}
	}
}