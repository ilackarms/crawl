package main

import (
	"encoding/json"
	"net"
	"log"
	"bufio"
	"fmt"
)

type PlayerData struct {
	Name  string `json:"Name"`
	X     int `json:"X"`
	Y     int `json:"Y"`
	Level int `json:"Level"`
}

func main() {
	pd := PlayerData{
		Name: "mircol",
		X: 1,
		Y: 2,
		Level: 3,
	}
	data, err := json.Marshal(pd)
	if err != nil {
		log.Fatal(err)
	}
	data = append([]byte("abcd\n"), data...)
	data = append(data, byte('\n'))
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func(){
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err == nil {
				//log.Printf("failed reading response %v", err)
				fmt.Printf("%s\n", line)
			}
		}
	}()
	if _, err := conn.Write(data); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote data: %s", data)
	
	select {}
}