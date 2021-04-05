package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
)

type GrpcServer struct {
	db db
}

func (g GrpcServer) GetTrainingBalance(ctx context.Context, request *users.GetTrainingBalanceRequest) (*users.GetTrainingBalanceResponse, error) {
	user, err := g.db.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.GetTrainingBalanceResponse{Amount: int64(user.Balance)}, nil
}

func (g GrpcServer) UpdateTrainingBalance(
	ctx context.Context,
	req *users.UpdateTrainingBalanceRequest,
) (*empty.Empty, error) {
	err := g.db.UpdateBalance(ctx, req.UserId, int(req.AmountChange))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update balance: %s", err))
	}

	return &empty.Empty{}, nil
}
