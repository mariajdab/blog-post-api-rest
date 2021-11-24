package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mariajdab/post-api-rest/api/models"
)

type mockMongoDB struct {
	data map[string]interface{}
}

func newMockMongoDB() *mockMongoDB {
	db := new(mockMongoDB)

	db.data = make(map[string]interface{})

	return db
}

func (m *mockMongoDB) CreatePost(value *models.Post) (interface{}, error) {
	id := uuid.New().String()
	m.data[id] = *value
	return id, nil
}

func (m *mockMongoDB) ReadPost(id string) (models.Post, error) {
	post, ok := m.data[id].(models.Post)
	if !ok {
		return models.Post{}, fmt.Errorf("not found")
	}

	return post, nil
}

func (m *mockMongoDB) ReadAllPosts() ([]models.Post, error) {
	var posts []models.Post
	for _, value := range m.data {
		post, ok := value.(models.Post)
		if ok {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (m *mockMongoDB) UpdatePost(idStr string, data map[string]interface{}) (int64, error) {
	return 0, nil
}

func (m *mockMongoDB) DeletePost(id string) (int64, error) {
	delete(m.data, id)

	if m.data[id] != nil {
		return -1, fmt.Errorf("could not remove from DB")
	}

	return 1, nil
}
