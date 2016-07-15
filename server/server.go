package main

import (
	"net"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"bufio"
	"encoding/json"
	"strings"
	tl "github.com/ilackarms/termloop"
	"time"
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

func main() {
	clients = make(map[string]*Client)
	levels = make(map[string]tl.Level)

	l, err := listen()
	if err != nil {
		panic(err)
	}
	log.Printf("listening on :9000")
	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})
	levels["0"] = level
	game.Screen().SetLevel(level)
	go func() {
		game.StartServerMode()
	}()
	go func() {
		for i := 0; i < 20; i++ {
			level.AddEntity(tl.NewRectangle(10 * i, 10 * i, 10 * i, 10 * i, tl.ColorBlue))
			time.Sleep(time.Second)
		}
		log.Printf("current entities: %v", level.Entities)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: failed to accept connection: %v", err)
			continue
		}
		go func() {
			go func(){
				//for {
					conn.Write([]byte("hihih\n"))
				//}
			}()
			defer conn.Close()
			log.Printf("new client connected: %v", conn)
			reader := bufio.NewReader(conn)
			clientID, err := reader.ReadString('\n')
			if err != nil {
				conn.Write([]byte("could not read client id: " + err.Error()))
				log.Printf("ERROR: could not read client id: " + err.Error())
				return
			}
			clientID = strings.TrimSuffix(clientID, "\n")
			log.Printf("new client %v connected", clientID)
			client, ok := clients[clientID]
			if !ok {
				//if client does not exist, expect playerinfo to be sent as next line
				data, err := reader.ReadBytes('\n')
				if err != nil {
					conn.Write([]byte("could not read playerdata: " + err.Error()))
					log.Printf("ERROR: could not read playerdata: " + err.Error())
					return
				}
				log.Printf("creating new player from data %s", data)
				var playerData PlayerData
				if err := json.Unmarshal(data, &playerData); err != nil {
					conn.Write([]byte("could not unmarshal playerdata: " + err.Error()))
					log.Printf("ERROR: could not unmarshal playerdata: " + err.Error())
					return
				}
				client = &Client{
					ID: clientID,
					conn: conn,
					PlayerData: playerData,
				}
				log.Printf("created new client %+v", client)
			}
			log.Printf("sending %v entities", len(level.Entities))
			levelData, err := json.Marshal(levels)
			if err != nil {
				conn.Write([]byte("could not marshal leveldata: " + err.Error()))
				log.Printf("ERROR: could not marshal leveldata: " + err.Error())
				return
			}
			if _, err := conn.Write(levelData); err != nil {
				conn.Write([]byte("could not send leveldata: " + err.Error()))
				log.Printf("ERROR: could not send leveldata: " + err.Error())
				return
			}
			conn.Write([]byte("\n"))
			clients[clientID] = client
			log.Printf("current clients: %+v", clients)
		}()
	}
}

func listen() (net.Listener, error) {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		return nil, errors.New("starting listener", err)
	}
	return l, nil
}