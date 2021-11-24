package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mariajdab/post-api-rest/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMockServer() *server {
	s := &server{
		server: &http.Server{},
		router: mux.NewRouter(),
		db:     newMockMongoDB(),
	}
	s.setRoutes()
	return s
}

func TestCreatePost(t *testing.T) {
	s := newMockServer()

	insertTime := time.Now()
	post := models.Post{
		CreatedAt: insertTime,
		Body:      "Body Test Create",
		UserName:  "Venusita",
	}

	payload, err := json.Marshal(&post)

	req, err := http.NewRequest(http.MethodGet, "/post", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreatePost)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var id string
	err = json.NewDecoder(rr.Body).Decode(&id)
	if err != nil {
		t.Fatal("Could not decode post:", err)
	}

	respPost, err := s.db.ReadPost(id)
	if err != nil {
		t.Fatal(err)
	}

	if respPost.CreatedAt.UTC() != post.CreatedAt.UTC() || respPost.Body != post.Body ||
		respPost.UserName != post.UserName {
		t.Errorf("db returned unexpected body: got %v want %v",
			respPost, post)
	}
}

func TestGetPost(t *testing.T) {
	s := newMockServer()

	insertTime := time.Now()
	post := models.Post{
		CreatedAt: insertTime,
		Body:      "Body Test Get",
		UserName:  "Venusita",
	}

	idRaw, err := s.db.CreatePost(&post)
	if err != nil {
		t.Fatal(err)
	}

	id, ok := idRaw.(string)
	if !ok {
		t.Fatal("Could not read ID")
	}

	req, err := http.NewRequest(http.MethodGet, "/post", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetPost)

	vars := map[string]string{
		"ID": id,
	}

	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var respPost models.Post
	err = json.NewDecoder(rr.Body).Decode(&respPost)
	if err != nil {
		t.Fatal("Could not decode post:", err)
	}

	if respPost.CreatedAt.UTC() != post.CreatedAt.UTC() || respPost.Body != post.Body ||
		respPost.UserName != post.UserName {
		t.Errorf("handler returned unexpected body: got %v want %v",
			respPost, post)
	}
}

func TestGetPosts(t *testing.T) {
	s := newMockServer()

	insertTime := time.Now()
	post := models.Post{
		CreatedAt: insertTime,
		Body:      "Body 1",
		UserName:  "Venusita",
	}

	_, err := s.db.CreatePost(&post)
	if err != nil {
		t.Fatal(err)
	}

	insertTime = time.Now()
	post = models.Post{
		CreatedAt: insertTime,
		Body:      "Body 2",
		UserName:  "Venusita",
	}

	_, err = s.db.CreatePost(&post)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetPosts)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var respPosts []models.Post
	err = json.NewDecoder(rr.Body).Decode(&respPosts)
	if err != nil {
		t.Fatal("Could not decode post:", err)
	}

	if len(respPosts) != 2 {
		t.Fatal("Unexpected len of slice:", err)
	}
}

func TestDeletePost(t *testing.T) {
	s := newMockServer()

	insertTime := time.Now()
	post := models.Post{
		CreatedAt: insertTime,
		Body:      "Post will be deleted",
		UserName:  "Venusita",
	}

	idRaw, err := s.db.CreatePost(&post)
	if err != nil {
		t.Fatal(err)
	}

	id := idRaw.(string)

	req, err := http.NewRequest(http.MethodDelete, "/post", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"ID": id,
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeletePost)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var deleteCount int
	err = json.NewDecoder(rr.Body).Decode(&deleteCount)
	if err != nil {
		t.Errorf("invalid response: %v", err)
	}

	if deleteCount != 1 {
		t.Errorf("invalid response body: got %v, expected 1", deleteCount)
	}

	_, err = s.db.ReadPost(id)
	if err.Error() != "not found" {
		t.Errorf("could not remove entry")
	}
}
