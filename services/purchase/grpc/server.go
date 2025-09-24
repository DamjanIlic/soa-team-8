package grpc

import (
	"context"
	"fmt"
	"purchase/pb"
	"purchase/service"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type PurchaseGRPCServer struct {
	pb.UnimplementedPurchaseServiceServer
	TokenService *service.TokenService
	CartService  *service.CartService
}

// Funkcija za dobijanje user ID iz JWT metadata
func (s *PurchaseGRPCServer) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("no metadata found")
	}

	auth := md.Get("authorization")
	if len(auth) == 0 {
		return uuid.Nil, fmt.Errorf("no authorization header")
	}

	tokenString := strings.Replace(auth[0], "Bearer ", "", 1)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id not found in token")
	}

	return uuid.Parse(userIDStr)
}

func (s *PurchaseGRPCServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Fallback: koristi user_id iz request-a ili JWT-a
	var userID uuid.UUID
	var err error

	if req.UserId != "" {
		userID, err = uuid.Parse(req.UserId)
	} else {
		userID, err = s.getUserIDFromContext(ctx)
	}

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

func (s *PurchaseGRPCServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	// Fallback: koristi user_id iz request-a ili JWT-a
	var userID uuid.UUID
	var err error

	if req.UserId != "" {
		userID, err = uuid.Parse(req.UserId)
	} else {
		userID, err = s.getUserIDFromContext(ctx)
	}

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
