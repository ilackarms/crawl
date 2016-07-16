package server

import (
	"net"
	tl "github.com/ilackarms/termloop"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
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
	game := startGame()

	//start multiplayer server
	if err := startServer(); err != nil {
		return errors.New("starting server", err)
	}

}

func startGame() *tl.Game {
	game := tl.NewGame()
	go game.StartServerMode()
	return game
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