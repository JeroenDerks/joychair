package joychair

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"gopkg.in/redis.v5"
)

type JoyServer struct {
	conn *redis.pubsub
}

type JoyNetEvent struct {
	x, y int8
}

func InitJoyServer() JoyServer {
	log.Printf("JoyServer with address: erhm yeah hardcoded...")

	client := redis.NewClient(&redis.Options{
		//Network:  "unix",
		//Addr:     "/tmp/redis.sock",
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pubsub, err := client.Subscribe("test")
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()

	server := JoyServer{conn: pubsub}

	return server
}

func (j *JoyServer) readLoop(c chan JoyNetEvent) {

	input := make([]byte, 2, 2)

	for {

		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Channel, msg.Payload)
		byteReader := bytes.NewReader(msg.Payload)

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
