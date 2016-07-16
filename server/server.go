package server

import (
	"net"
	tl "github.com/ilackarms/termloop"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"github.com/ilackarms/crawl/game"
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

var clients map[string]*Client
var levels map[string]tl.Level

func init() {
	clients = make(map[string]*Client)
	levels = make(map[string]tl.Level)
}

func Start() error {
	//start game
	g := startGame()

	//start multiplayer server
	if err := startServer(); err != nil {
		return errors.New("starting server", err)
	}

	syncLevel := func(level game.Level) {
		for _, client := range clients {
			if err := func() error {
				data, err := game.SerializeLevel(level)
				if err != nil {
					return errors.New("serializing level", err)
				}
				if err := client.conn.Write(data); err != nil {
					return errors.New("writing level data to client", err)
				}
				return nil
			}(); err != nil {
				log.Printf("ERROR: failed syncing level with client: %v", err)
			}
		}
	}
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
	go func(){
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