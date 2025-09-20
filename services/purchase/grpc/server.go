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
