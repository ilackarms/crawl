package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/emc-advanced-dev/pkg/errors"
	"log"
	"net"
)

//Generate a Crawl TCP Protocol Packet from json object
func generatePacket(message interface{}, messageType byte) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("message was not json serializable", err)
	}
	size := uint32(len(data))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, size)
	packet := append(append(bs, messageType), data...)
	if false {
		log.Printf("packet: %s", packet)
	}
	return packet, nil
}

//Send Message on TCP Connection
func SendMessage(conn net.Conn, message interface{}, messageType byte) error {
	packet, err := generatePacket(message, messageType)
	if err != nil {
		return errors.New("generating packet", err)
	}
	if _, err := conn.Write(packet); err != nil {
		return errors.New("sending packet", err)
	}
	return nil
}

//Read Message on TCP Connection. Returns raw json string, message type, and error or nil
func ReadMessage(conn net.Conn) ([]byte, byte, error) {
	bs := make([]byte, 4)
	if _, err := conn.Read(bs); err != nil {
		return nil, byte(0), errors.New("packet size from connection", err)
	}
	messageType := make([]byte, 1)
	if _, err := conn.Read(messageType); err != nil {
		return nil, byte(0), errors.New("packet type from connection", err)
	}
	size := binary.LittleEndian.Uint32(bs)
	message := make([]byte, size)
	if _, err := conn.Read(message); err != nil {
		return nil, byte(0), errors.New(fmt.Sprintf("reading message of size %v from connection", size), err)
	}
	return message, messageType[0], nil
}
