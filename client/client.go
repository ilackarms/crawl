package client

import (
	"encoding/json"
	"fmt"
	"github.com/ilackarms/crawl/game"
	"github.com/ilackarms/crawl/protocol"
	tl "github.com/ilackarms/termloop"
	"log"
	"net"
)

var player *game.Player

func Start(name, serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed connecting to server: %v", err)
	}
	player = game.NewPlayer(name, tl.NewEntity(1, 1, 1, 1), conn)
	player.SetDescription("this is a description about some things.\nfurthermore, newlines don't work at all. isn't that a funny happenstance coincedence?")
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
				continue
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
					screenWidth, screenHeight := g.Screen().Size()
					x, y := player.Position()
					level.SetOffset(screenWidth/2-x, screenHeight/2-y)
					player.InputText.SetPosition(x-len(player.InputText.GetText())/2, y-1+screenHeight/2)
					level.RemoveEntity(playerRep)
					level.AddEntity(player)
					if description, ok := level.Descriptions[fmt.Sprintf("%d,%d", x, y)]; ok {
						player.SetDescription(description)
					}
					playerFound = true
					break
				}
			}
			if !playerFound {
				log.Printf("level received: %v", string(message))
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
