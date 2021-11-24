package database

import (
	"context"
	"fmt"
	"github.com/mariajdab/post-api-rest/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWrapper interface {
	CreatePost(*models.Post) (interface{}, error)
	ReadPost(string) (models.Post, error)
	ReadAllPosts() ([]models.Post, error)
	UpdatePost(string, map[string]interface{}) (int64, error)
	DeletePost(string) (int64, error)
}

type Mongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoClient() MongoWrapper {
	URL := fmt.Sprintf(os.Getenv("MONGODB_URI"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		panic(err)
	}

	collection := client.Database("db-mongo").Collection("blog-posts")

	return &Mongo{
		Client:     client,
		Collection: collection,
	}
}

func (m *Mongo) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.Client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func (m *Mongo) CreatePost(value *models.Post) (interface{}, error) {
	ctx := context.Background()
	collection := m.Collection

	result, err := collection.InsertOne(ctx, value)
	if err != nil {
		return models.Post{}, err
	}

	return result.InsertedID, nil
}

func (m *Mongo) ReadPost(idStr string) (models.Post, error) {
	var res models.Post
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return models.Post{}, fmt.Errorf("could not parse Id:" + err.Error())
	}

	ctx := context.Background()
	postsCollection := m.Collection
	err = postsCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&res)
	if err != nil {
		return models.Post{}, err
	}

	return res, nil
}

func (m *Mongo) ReadAllPosts() ([]models.Post, error) {
	var res []models.Post
	ctx := context.Background()
	postsCollection := m.Collection

	cursor, err := postsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Mongo) UpdatePost(idStr string, data map[string]interface{}) (int64, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return -1, fmt.Errorf("could not parse Id:" + err.Error())
	}

	filter := bson.D{{"_id", id}}

	postsCollection := m.Collection

	update := bson.D{{"$set", bson.D{{"body", data["body"]}}}}
	updatePost, err := postsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return -1, fmt.Errorf("could not update record: %v", err)
	}

	return updatePost.ModifiedCount, nil
}

func (m *Mongo) DeletePost(idStr string) (int64, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return -1, fmt.Errorf("could not parse Id:" + err.Error())
	}

	postsCollection := m.Collection
	opts := options.Delete().SetCollation(&options.Collation{})
	result, err := postsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}}, opts)
	if err != nil {
		return -1, fmt.Errorf("an error happened: %v", err.Error())
	}

	return result.DeletedCount, nil
}
