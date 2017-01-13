package main

import (
	"github.com/BurntSushi/toml"
	"github.com/JeroenDerks/joychair"
	"github.com/tarm/serial"
	"log"
	"os"
)

type Config struct {
	Chair serial.Config
}

func main() {

	if len(os.Args) <= 1 {
		log.Fatal("Provide configuration path as first argument")
	}

	configPath := os.Args[1]

	config := Config{}

	log.Printf("Reading configuration from %s", configPath)

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(err)
	}

	chair := joychair.InitChair(&config.Chair)

	chair.Loop()

	log.Printf("Bye")

}

func redisSubscribeLoop(c chan joychair.JoyNetEvent) { // reads messages from Redis, checks for priority, and sends to channel

	event := new(joychair.JoyNetEvent)

	//something here
	c <- *event

}
