package repo

import (
	"blog/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepository struct {
	Collection *mongo.Collection
}

// Kreiranje novog like-a
func (r *LikeRepository) Create(like *model.Like) error {
	_, err := r.Collection.InsertOne(context.TODO(), like)
	return err
}

// Brisanje like-a po userID i blogID
func (r *LikeRepository) Delete(userID, blogID string) error {
	_, err := r.Collection.DeleteOne(context.TODO(), bson.M{"user_id": userID, "blog_id": blogID})
	return err
}

// Provera da li like već postoji
func (r *LikeRepository) Exists(userID, blogID string) (bool, error) {
	count, err := r.Collection.CountDocuments(context.TODO(), bson.M{"user_id": userID, "blog_id": blogID})
	return count > 0, err
}

// Broj like-ova po blogID-u
func (r *LikeRepository) CountByBlogID(blogID string) (int64, error) {
	return r.Collection.CountDocuments(context.TODO(), bson.M{"blog_id": blogID})
}
