package server

import (
	"fmt"
	"github.com/emc-advanced-dev/pkg/errors"
	"github.com/ilackarms/crawl/game"
	tl "github.com/ilackarms/termloop"
	"log"
	"net"
)

var world *game.World

func init() {
	world = game.NewWorld()
}

func Start() {
	//start game
	world.StartGame()

	//test - create  & set the current level
	level1 := game.NewLevel(tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	}))
	level2 := game.NewDungeonLevel(20, 50)
	level1.AddEntity(tl.NewRectangle(20, -20, 30, 30, tl.ColorBlue))
	level1.AddEntity(game.NewTrigger(10, -10,
		map[game.Position]func(player *game.PlayerRep){
			game.Position{X: 3, Y: 0}: func(player *game.PlayerRep) {
				for _, c := range world.Clients {
					player := c.PlayerRep
					levelUUID := level2.UUID
					player.W.Levels[levelUUID].AddEntity(player)
					player.Iq.Push(game.InputMessage{
						CustomEvent: func() {
							//center player
							player.PrevX = 0
							player.PrevY = 0
							player.Entity.SetPosition(0, 0)
						},
					})
					player.W.SetLevel(levelUUID)
				}
			},
		},
		game.DungeonEntrance, tl.ColorWhite))
	level1.AddEntity(game.NewTriggerArea(10, -30, 40, 40, func(player *game.PlayerRep) {
		for _, c := range world.Clients {
			player := c.PlayerRep
			levelUUID := level2.UUID
			player.W.Levels[levelUUID].AddEntity(player)
			player.Iq.Push(game.InputMessage{
				CustomEvent: func() {
					//center player
					player.PrevX = 0
					player.PrevY = 0
					player.Entity.SetPosition(0, 0)
				},
			})
			player.W.SetLevel(levelUUID)
		}
	}))
	level1.Descriptions[fmt.Sprintf("%d,%d", 5, 5)] = "try this"
	level1.Descriptions[fmt.Sprintf("%d,%d", -5, -5)] = "try that"

	world.AddLevel(level1)
	world.AddLevel(level2)
	world.SetLevel(level1.UUID)
	//test

	//accept client connections
	acceptConections()

}

func acceptConections() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(errors.New("starting listener", err))
	}
	log.Printf("listening on :9000")
	func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatalf("ERROR: failed to accept connection: %v", err)
			}
			go world.HandleClient(conn)
		}
	}()
}
