package server

import (
	"net"
	tl "github.com/ilackarms/termloop"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"github.com/ilackarms/crawl/game"
	"github.com/ilackarms/crawl/protocol"
)

type Client struct {
	ID         string `json:"ID"`
	PlayerData PlayerData `json:"PlayerData"`
	conn       net.Conn `json:"-"`
}

type PlayerData struct {
	Name  string `json:"Name"`
	X     int `json:"X"`
	Y     int `json:"Y"`
	Level int `json:"Level"`
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

}