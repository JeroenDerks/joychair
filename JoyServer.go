package joychair

import (
	"log"
	"net"
)

type JoyServer struct {
	conn *net.UDPConn
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

func (j *JoyServer) send(data *ChairResponse) {
	//lets send some chair data!
	j.conn.Write(data.bytes())
}
