package client

import (
	tl "github.com/ilackarms/termloop"
	"github.com/ilackarms/crawl/game"
	"net"
	"log"
	"github.com/ilackarms/crawl/protocol"
	"encoding/json"
)

var player *game.Player

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
	if err := protocol.SendMessage(conn, login, game.Login.GetByte()); err != nil {
		log.Fatalf("failed logging in: %v", err)
	}
	g := tl.NewGame()
	go g.Start()
	for {
		message, messageType, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Fatalf("failed to read message from server: %v", err)
		}
		switch messageType {
		case game.LevelUpdate.GetByte():
			var levelUpdate game.LevelChangeMessage
			if err := json.Unmarshal(message, &levelUpdate); err != nil {
				log.Fatalf("ERROR: unmarshalling level update: %v", err)
			}
			level, err := game.DeserializeLevel(levelUpdate.LevelData)
			if err != nil {
				log.Fatalf("ERROR: deserializing level: %v", err)
			}
			var playerFound bool
			for _, entity := range level.Entities {
				//swap out player rep for player on local end
				if entity.GetUUID() == player.GetUUID() {
					playerRep, ok := entity.(*game.PlayerRep)
					if !ok {
						log.Fatalf("ERROR: same uuid but not a player. what?")
					}
					player.SetPosition(playerRep.Position())
					player.SetLevel(level)
					level.RemoveEntity(playerRep)
					level.AddEntity(player)
					playerFound = true
					break
				}
			}
			if !playerFound {
				log.Fatalf("ERROR: player not found in level!")
			}
			g.Screen().SetLevel(level)
		}
		/*
		TODO: go to server, make sure server replies to login with the first level
		so the client can start drawing right away and the input cycle can begin
		*/
	}
}
