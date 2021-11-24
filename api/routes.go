package api

import (
	"net/http"
)

func (s *server) setRoutes() {
	s.router.HandleFunc("/post", s.CreatePost).Methods(http.MethodPost)
	s.router.HandleFunc("/post/{ID}", s.GetPost).Methods(http.MethodGet)
	s.router.HandleFunc("/posts", s.GetPosts).Methods(http.MethodGet)
	s.router.HandleFunc("/post/{ID}", s.UpdatePost).Methods(http.MethodPut)
	s.router.HandleFunc("/post/{ID}", s.DeletePost).Methods(http.MethodDelete)
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello"))
	})
}
