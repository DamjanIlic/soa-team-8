package repo

import (
	"blog/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository struct {
	Collection *mongo.Collection
}

// Kreiranje novog komentara
func (r *CommentRepository) Create(comment *model.Comment) error {
	_, err := r.Collection.InsertOne(context.TODO(), comment)
	return err
}

// Dohvatanje komentara po blogID-u (string)
func (r *CommentRepository) GetByBlogID(blogID string) ([]model.Comment, error) {
	cursor, err := r.Collection.Find(context.TODO(), bson.M{"blog_id": blogID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var comments []model.Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}
