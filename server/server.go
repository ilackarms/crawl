package server

import (
	"net"
	tl "github.com/ilackarms/termloop"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"github.com/ilackarms/crawl/game"
	"github.com/ilackarms/crawl/game/objects"
)

var world *game.World

func init() {
	world = game.NewWorld()
}

func Start() {
	//start game
	world.StartGame()

	//test - create  & set the current level
	level1 := &game.Level{
		BaseLevel: tl.NewBaseLevel(tl.Cell{
			Bg: tl.ColorGreen,
			Fg: tl.ColorBlack,
			Ch: 'v',
		}),
	}
	level2 := &game.Level{
		BaseLevel: tl.NewBaseLevel(tl.Cell{
			Bg: tl.RgbTo256Color(20,20,20),
			Fg: tl.RgbTo256Color(110,110,110),
			Ch: '░',
		}),
	}
	level1.AddEntity(tl.NewRectangle(20, 20, 30, 30, tl.ColorBlue))
	level1.AddEntity(objects.NewDungeonEntrance(10, 10, tl.ColorWhite, level2.UUID))

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
