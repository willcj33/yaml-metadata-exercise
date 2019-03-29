package main

import (
	"fmt"
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/handlers"
)

func main() {
	cfg, _ := config.GetConfig()

	metadataServer := handlers.NewMetadataServer(cfg)
	mux := metadataServer.CreateMux()
	httpListenAddr := fmt.Sprintf("%s:%s", cfg.HTTPListenHost, cfg.HTTPListenPort)

	// Start listening
	server := http.Server{
		Addr:    httpListenAddr,
		Handler: mux,
	}

	fmt.Printf("%s listening on %s", cfg.ServerName, httpListenAddr)
	server.ListenAndServe()
}
