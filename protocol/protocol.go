package protocol

import (
	"encoding/json"
	"encoding/binary"
	"github.com/emc-advanced-dev/pkg/errors"
	"net"
	"log"
	"fmt"
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
	return append(bs, messageType, data...), nil
}

//Send Message on TCP Connection
func SendMessage(conn net.Conn, message interface{}, messageType byte) error {
	packet, err := generatePacket(message, messageType)
	if err != nil {
		return errors.New("generating packet", err)
	}
	n, err := conn.Write(packet)
	if err != nil {
		return errors.New("sending packet", err)
	}
	log.Printf("%v bytes sent to %v", n, conn.RemoteAddr().String())
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