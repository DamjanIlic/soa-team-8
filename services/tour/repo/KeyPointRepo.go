package repo

import (
	"tour/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeyPointRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *KeyPointRepository) Create(keyPoint *model.KeyPoint) error {
	return r.DatabaseConnection.Create(keyPoint).Error
}

func (r *KeyPointRepository) GetByTourID(tourID uuid.UUID) ([]model.KeyPoint, error) {
	var keyPoints []model.KeyPoint

	if err := r.DatabaseConnection.Where("tour_id = ?", tourID).Order("\"order\" ASC").Find(&keyPoints).Error; err != nil {
		return nil, err
	}
	return keyPoints, nil
}

func (r *KeyPointRepository) GetByID(id string) (*model.KeyPoint, error) {
	var keyPoint model.KeyPoint

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if err := r.DatabaseConnection.First(&keyPoint, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &keyPoint, nil
}

func (r *KeyPointRepository) Update(keyPoint *model.KeyPoint) error {
	return r.DatabaseConnection.Save(keyPoint).Error
}

func (r *KeyPointRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.DatabaseConnection.Delete(&model.KeyPoint{}, "id = ?", uid).Error
}

func (r *KeyPointRepository) DeleteByTourID(tourID uuid.UUID) error {
	return r.DatabaseConnection.Where("tour_id = ?", tourID).Delete(&model.KeyPoint{}).Error
}

// pronalazi najveci order za turu
func (r *KeyPointRepository) GetMaxOrder(tourID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.DatabaseConnection.Model(&model.KeyPoint{}).
		Where("tour_id = ?", tourID).
		Select("COALESCE(MAX(\"order\"), 0)"). // ako nema tačaka, vraća 0
		Scan(&maxOrder).Error
	if err != nil {
		return 0, err
	}
	return maxOrder, nil
}
