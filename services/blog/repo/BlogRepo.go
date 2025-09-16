package repo

import (
	"blog/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	Collection *mongo.Collection
}

// Kreiranje novog bloga
func (r *BlogRepository) Create(blog *model.Blog) error {
	_, err := r.Collection.InsertOne(context.TODO(), blog)
	return err
}

// Dohvatanje svih blogova
func (r *BlogRepository) GetAll() ([]model.Blog, error) {
	cursor, err := r.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []model.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

// Dohvatanje bloga po string ID-u
func (r *BlogRepository) Get(id string) (*model.Blog, error) {
	var blog model.Blog
	err := r.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&blog)
	return &blog, err
}

// Update bloga po string ID-u
func (r *BlogRepository) Update(blog *model.Blog) error {
	_, err := r.Collection.ReplaceOne(
		context.TODO(),
		bson.M{"_id": blog.ID},
		blog,
	)
	return err
}
