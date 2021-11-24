package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mariajdab/post-api-rest/api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (s *server) CreatePost(w http.ResponseWriter, r *http.Request) {
	var data models.Post

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not decode from JSON:"+err.Error())
		return
	}

	data.CreatedAt = time.Now()
	insertedId, err := s.db.CreatePost(&data)
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not insert into DB:"+err.Error())
		return
	}

	s.SuccessResponse(w, http.StatusOK, insertedId)
}

func (s *server) GetPost(w http.ResponseWriter, r *http.Request) {
	var data models.Post
	id := mux.Vars(r)["ID"]

	data, err := s.db.ReadPost(id)
	if err == mongo.ErrNoDocuments {
		s.SuccessResponse(w, http.StatusNotFound, nil)
		return
	} else if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not read from DB:"+err.Error())
		return
	}

	s.SuccessResponse(w, http.StatusOK, data)
}

func (s *server) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	posts, err := s.db.ReadAllPosts()
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not access DB: "+err.Error())
		return
	}

	if posts == nil {
		s.SuccessResponse(w, http.StatusNotFound, nil)
		return
	}

	s.SuccessResponse(w, http.StatusOK, posts)
	return
}

func (s *server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	id := mux.Vars(r)["ID"]

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not read the request:"+err.Error())
		return
	}

	result, err := s.db.UpdatePost(id, data)
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "Could not updated entry:"+err.Error())
		return
	}

	s.SuccessResponse(w, http.StatusOK, result)
	return
}

func (s *server) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["ID"]

	result, err := s.db.DeletePost(id)
	if err != nil {
		s.processError(w, http.StatusInternalServerError, "could not remove record: "+err.Error())
	}

	if result == 0 {
		s.SuccessResponse(w, http.StatusNotFound, nil)
		return
	}

	s.SuccessResponse(w, http.StatusOK, result)
	return
}
