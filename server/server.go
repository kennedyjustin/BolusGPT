package server

import (
	"net/http"
	"sync"

	"github.com/kennedyjustin/BolusGPT/dexcom"
	"github.com/kennedyjustin/BolusGPT/jsonfile"
)

type Server struct {
	mu           sync.Mutex
	server       *http.Server
	db           *jsonfile.JSONFile[Me]
	dexcomClient *dexcom.Client
}

type ServerInput struct {
	FilePath       string
	DexcomUsername string
	DexcomPassword string
}

func NewServer(input ServerInput) (*Server, error) {
	server := &Server{}

	db, err := jsonfile.LoadOrNew[Me](input.FilePath)
	if err != nil {
		return nil, err
	}
	server.db = db

	dexcomClient, err := dexcom.NewClient(dexcom.ClientInput{
		Username: input.DexcomUsername,
		Password: input.DexcomPassword,
	})
	if err != nil {
		return nil, err
	}
	server.dexcomClient = dexcomClient

	mux := http.NewServeMux()
	mux.HandleFunc("GET /me", server.MeHandlerGet)
	mux.HandleFunc("PATCH /me", server.MeHandlerPatch)
	mux.HandleFunc("POST /dose", server.DoseHandler)
	httpServer := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.server = httpServer

	return server, nil
}

func (s *Server) Start() {
	s.server.ListenAndServe()
}
