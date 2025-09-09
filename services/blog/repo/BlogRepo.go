package repo

import (
	"blog/model"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *BlogRepository) Create(blog *model.Blog) error {
	return r.DatabaseConnection.Create(blog).Error
}

func (r *BlogRepository) GetAll() ([]model.Blog, error) {
	var blogs []model.Blog
	result := r.DatabaseConnection.Find(&blogs)
	return blogs, result.Error
}

func (r *BlogRepository) Get(id string) (*model.Blog, error) {
	var blog model.Blog
	result := r.DatabaseConnection.First(&blog, "id = ?", id)
	return &blog, result.Error
}

func (r *BlogRepository) Update(blog *model.Blog) error {
	return r.DatabaseConnection.Save(blog).Error
}
