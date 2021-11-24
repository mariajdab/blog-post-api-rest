package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mariajdab/post-api-rest/internal/database"
	"log"
	"net/http"
)

type server struct {
	server *http.Server
	router *mux.Router
	db     database.MongoWrapper
}

func NewServer() *server {
	s := &server{
		server: &http.Server{},
		router: mux.NewRouter(),
		db:     database.NewMongoClient(),
	}

	s.setRoutes()
	s.server.Addr = ":8080"
	s.server.Handler = s.router

	return s
}

func (s *server) Run() {
	log.Println("Running server on: ", s.server.Addr)
	log.Fatal(s.server.ListenAndServe())
}

func (s *server) SuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		panic(err)
	}
}

func (s *server) processError(w http.ResponseWriter, code int, message string) {
	log.Println("[ERROR]\t", message)
	s.SuccessResponse(w, code, message)
}
