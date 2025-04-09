package main

import (
	"log"
	"os"

	"github.com/kennedyjustin/BolusGPT/server"
)

const Filepath = "me.json"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s, err := server.NewServer(server.ServerInput{
		FilePath:       Filepath,
		DexcomUsername: os.Getenv("DEXCOM_USERNAME"),
		DexcomPassword: os.Getenv("DEXCOM_PASSWORD"),
		BearerToken:    os.Getenv("BEARER_TOKEN"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
}
