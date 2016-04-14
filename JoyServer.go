package joychair

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type JoyServer struct {
	conn *net.UDPConn
}

type JoyNetEvent struct {
	x, y int8
}

func InitJoyServer() JoyServer {
	log.Printf("JoyServer with address: erhm yeah hardcoded...")

	ServerAddr, err := net.ResolveUDPAddr("udp", "195.88.37.61:3001")
	if err != nil {
		log.Println("S Error: ", err)
	}

	conn, err := net.DialUDP("udp", nil, ServerAddr)
	if err != nil {
		log.Println("D Error: ", err)
	}

	server := JoyServer{conn: conn}

	return server
}

func (j *JoyServer) ReadLoop(c chan JoyNetEvent) {

	input := make([]byte, 2, 2)

	if j.conn == nil {
		return
	}

	for {

		nBytes, err := j.conn.Read(input)

		if nBytes != 2 {
			continue
		}

		byteReader := bytes.NewReader(input)

		event := new(JoyNetEvent)

		binary.Read(byteReader, binary.LittleEndian, &event.x)

		err = binary.Read(byteReader, binary.LittleEndian, &event.y)

		if err != nil {
			log.Fatal("binary.Read failed:", err)
		}

		log.Print("I has a net event: %+v", event)

		c <- *event
	}
}

func (j *JoyServer) send(data *ChairResponse) {
	//lets send some chair data!
	go func(b []byte) {
		j.conn.Write(b)
	}(data.bytes())
}
