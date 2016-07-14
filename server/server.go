package main

import (
	"net"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"bufio"
	"encoding/json"
	"strings"
)

type Client struct {
	ID string `json:"ID"`
	PlayerData PlayerData `json:"PlayerData"`
	conn net.Conn `json:"-"`
}

type PlayerData struct {
	Name  string `json:"Name"`
	X     int `json:"X"`
	Y     int `json:"Y"`
	Level int `json:"Level"`
}

var clients map[string]*Client

func main() {
	clients = make(map[string]*Client)
	l, err := listen()
	if err != nil {
		panic(err)
	}
	log.Printf("listening on :9000")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: failed to accept connection: %v", err)
			continue
		}
		defer conn.Close()
		log.Printf("new client connected: %v", conn)
		reader := bufio.NewReader(conn)
		clientID, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("could not read client id: "+err.Error()))
			conn.Close()
			log.Printf("ERROR: could not read client id: "+err.Error())
			continue
		}
		clientID = strings.TrimSuffix(clientID, "\n")
		log.Printf("new client %v connected", clientID)
		client, ok := clients[clientID]
		if !ok {
			//if client does not exist, expect playerinfo to be sent as next line
			data, err := reader.ReadBytes('\n')
			if err != nil {
				conn.Write([]byte("could not read playerdata: "+err.Error()))
				conn.Close()
				log.Printf("ERROR: could not read playerdata: "+err.Error())
				continue
			}
			log.Printf("creating new player from data %s", data)
			var playerData PlayerData
			if err := json.Unmarshal(data, &playerData); err != nil {
				conn.Write([]byte("could not unmarshal playerdata: "+err.Error()))
				conn.Close()
				log.Printf("ERROR: could not unmarshal playerdata: "+err.Error())
				continue
			}
			client = &Client{
				ID: clientID,
				conn: conn,
				PlayerData: playerData,
			}
			log.Printf("created new client %+v", client)
		}
		clients[clientID] = client
		log.Printf("current clients: %+v", clients)
	}
}

func listen() (net.Listener, error) {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		return nil, errors.New("starting listener", err)
	}
	return l, nil
}