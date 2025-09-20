package repo

import (
	"purchase/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository struct {
	DB *gorm.DB
}

func (r *TokenRepository) Create(token *model.TourPurchaseToken) error {
	return r.DB.Create(token).Error
}

func (r *TokenRepository) GetByID(id uuid.UUID) (*model.TourPurchaseToken, error) {
	var token model.TourPurchaseToken
	if err := r.DB.First(&token, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) GetByTourist(touristID uuid.UUID) ([]model.TourPurchaseToken, error) {
	var tokens []model.TourPurchaseToken
	if err := r.DB.Where("tourist_id = ?", touristID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *TokenRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.TourPurchaseToken{}, "id = ?", id).Error
}
