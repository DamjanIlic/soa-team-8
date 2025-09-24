package grpc

import (
	"context"
	"purchase/pb" // import generirane proto fajlove
	"purchase/service"

	"github.com/google/uuid"
)

type PurchaseGRPCServer struct {
	pb.UnimplementedPurchaseServiceServer
	TokenService *service.TokenService
	CartService  *service.CartService // dodano za GetCart
}

func (s *PurchaseGRPCServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	tokens, err := s.TokenService.Checkout(userID)
	if err != nil {
		return nil, err
	}

	var grpcTokens []*pb.TourPurchaseToken
	for _, token := range tokens {
		grpcTokens = append(grpcTokens, &pb.TourPurchaseToken{
			TourId:    token.TourID.String(),
			TouristId: token.TouristID.String(),
			Token:     token.Token,
			CreatedAt: token.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &pb.CheckoutResponse{
		Tokens:  grpcTokens,
		Message: "checkout successful",
	}, nil
}

// Nova GetCart metoda
func (s *PurchaseGRPCServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	// Dohvati cart po tourist ID
	cart, err := s.CartService.GetByTouristID(userID)
	if err != nil {
		return nil, err
	}

	// Konvertuj cart items u protobuf format
	var grpcItems []*pb.CartItem
	for _, item := range cart.Items {
		grpcItems = append(grpcItems, &pb.CartItem{
			Id:     item.ID.String(),
			TourId: item.TourID.String(),
			Name:   item.Name,
			Price:  item.Price,
		})
	}

	// Izračunaj total
	total, err := s.CartService.GetTotal(cart.ID)
	if err != nil {
		return nil, err
	}

	return &pb.GetCartResponse{
		Id:        cart.ID.String(),
		TouristId: cart.TouristID.String(),
		Items:     grpcItems,
		Total:     total,
	}, nil
}
