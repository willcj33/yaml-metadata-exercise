package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/handlers"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Printf(dir)

	cfg, _ := config.GetConfig()

	metadataServer := handlers.NewMetadataServer(cfg)
	mux := metadataServer.CreateMux()
	httpListenAddr := fmt.Sprintf("%s:%s", cfg.HTTPListenHost, cfg.HTTPListenPort)

	// Start listening
	server := http.Server{
		Addr:    httpListenAddr,
		Handler: mux,
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.RemoveAll(fmt.Sprintf("%s/applicationMetadata.bleve", dir))
		os.Exit(1)
	}()
	fmt.Printf("%s listening on %s", cfg.ServerName, httpListenAddr)
	server.ListenAndServe()
}
