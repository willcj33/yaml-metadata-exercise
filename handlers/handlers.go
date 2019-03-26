package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/willcj33/yaml-metadata-exercise/db"

	"github.com/willcj33/yaml-metadata-exercise/config"
)

type httpHandlerWithStore = func(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, config config.Config)

//MetadataServer contains routes and config for establishing the server
type MetadataServer struct {
	routes map[string]map[string]http.Handler
	config config.Config
	store  *db.MetadataStore
}

func NewMetadataServer(cfg config.Config) *MetadataServer {
	server := &MetadataServer{
		routes: map[string]map[string]http.Handler{},
		config: cfg,
		store:  db.GetStore(cfg),
	}

	server.addRoute("/", http.MethodGet, func(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, cfg config.Config) {
		Health(w, r)
	})
	server.addRoute("/application/metadata", http.MethodPost, PostMetadata)
	server.addRoute("/application/metadata", http.MethodPut, PutMetadata)
	server.addRoute("/application/metadata", http.MethodGet, GetMetadata)

	return server
}

func (server *MetadataServer) CreateMux() http.Handler {
	serverMux := mux.NewRouter()
	for method, routes := range server.routes {
		for route, fn := range routes {
			serverMux.Handle(route, fn).Methods(method)
		}
	}

	return serverMux
}

func (server *MetadataServer) addRoute(route, method string, handlerFunc httpHandlerWithStore) {
	if _, ok := server.routes[method]; !ok {
		server.routes[method] = map[string]http.Handler{}
	}
	switch method {
	case http.MethodGet:
		server.routes[method][route] = httpGet(handlerFunc, server.store, server.config)
	case http.MethodPost:
		server.routes[method][route] = httpPost(handlerFunc, server.store, server.config)
	case http.MethodPut:
		server.routes[method][route] = httpPut(handlerFunc, server.store, server.config)
	default:
		fmt.Printf("Could not add route: %s:%s", method, route)
	}
}

func httpPut(next httpHandlerWithStore, store *db.MetadataStore, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r, store, cfg)
	})
}

func httpPost(next httpHandlerWithStore, store *db.MetadataStore, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r, store, cfg)
	})
}

func httpGet(next httpHandlerWithStore, store *db.MetadataStore, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next(w, r, store, cfg)
	})
}
