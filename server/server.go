package server

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/kennedyjustin/BolusGPT/dexcom"
	"github.com/kennedyjustin/BolusGPT/jsonfile"
)

type Server struct {
	mu             sync.Mutex
	server         *http.Server
	db             *jsonfile.JSONFile[Me]
	dexcomClient   *dexcom.Client
	bearerToken    string
	serverCertPath string
	serverkeyPath  string
}

type ServerInput struct {
	FilePath       string
	DexcomUsername string
	DexcomPassword string
	BearerToken    string
	ServerCertPath string
	ServerKeyPath  string
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
	mux.HandleFunc("GET /me", server.Auth(server.MeHandlerGet))
	mux.HandleFunc("PATCH /me", server.Auth(server.MeHandlerPatch))
	mux.HandleFunc("POST /dose", server.Auth(server.DoseHandler))
	httpServer := &http.Server{
		Handler: mux,
		Addr:    ":443",
	}
	server.server = httpServer

	server.bearerToken = input.BearerToken
	server.serverCertPath = input.ServerCertPath
	server.serverkeyPath = input.ServerKeyPath

	return server, nil
}

func (s *Server) Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		headerSlice := strings.Split(authHeader, "Bearer ")
		if authHeader == "" || len(headerSlice) != 2 || headerSlice[1] != s.bearerToken {
			print(headerSlice[1], s.bearerToken)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func (s *Server) Start() {
	err := s.server.ListenAndServeTLS(s.serverCertPath, s.serverkeyPath)
	if err != nil {
		log.Fatalln(err)
	}
}
