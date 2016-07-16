package client

import (
	tl "github.com/ilackarms/termloop"
	"github.com/ilackarms/crawl/game"
	"net"
	"log"
	"github.com/ilackarms/crawl/protocol"
)

var player game.Player

func Start(name, serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed connecting to server: %v", err)
	}
	player = game.NewPlayer(name, tl.NewEntity(1, 1, 1, 1), conn)
	//login
	login := game.LoginMessage{
		Name: name,
		UUID: player.GetUUID(),
	}
	if err := protocol.SendMessage(conn, login, game.Login); err != nil {
		log.Fatalf("failed logging in: %v", err)
	}
	g := tl.NewGame()
	go g.Start()

	for {
		/*
		TODO: go to server, make sure server replies to login with the first level
		so the client can start drawing right away and the input cycle can begin
		*/
	}
}
