package main

import (
	"log"
	"os"

	"github.com/kennedyjustin/BolusGPT/server"
)

const Filepath = "me.json"

func main() {
	_, err := server.NewServer(server.ServerInput{
		FilePath:       Filepath,
		DexcomUsername: os.Getenv("DEXCOM_USERNAME"),
		DexcomPassword: os.Getenv("DEXCOM_PASSWORD"),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
