package service

import (
	"purchase/model"
	"purchase/repo"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	TokenRepo *repo.TokenRepository
	CartRepo  *repo.CartRepository
	ItemRepo  *repo.ItemRepository
}

// checkout generise token za svaku stavku u korpi
func (s *TokenService) Checkout(touristID uuid.UUID) ([]model.TourPurchaseToken, error) {
	cart, err := s.CartRepo.GetByTouristID(touristID)
	if err != nil {
		return nil, err
	}

	var tokens []model.TourPurchaseToken
	for _, item := range cart.Items {
		token := model.TourPurchaseToken{
			TourID:    item.TourID,
			TouristID: touristID,
			Token:     uuid.New().String(),
			CreatedAt: time.Now(),
		}
		if err := s.TokenRepo.Create(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	// obrise sve stavke iz baze podataka
	for _, item := range cart.Items {
		if err := s.ItemRepo.Delete(item.ID); err != nil {
			return nil, err
		}
	}

	// isprazni korpu i postavi total na 0
	cart.Items = []model.OrderItem{}
	cart.Total = 0
	if err := s.CartRepo.Update(cart); err != nil {
		return nil, err
	}

	return tokens, nil
}

// vrati sve tokene koje turista ima
func (s *TokenService) GetTokensForTourist(touristID uuid.UUID) ([]model.TourPurchaseToken, error) {
	return s.TokenRepo.GetByTourist(touristID)
}
