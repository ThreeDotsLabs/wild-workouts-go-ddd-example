package adapters

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
)

type UsersGrpc struct {
	client users.UsersServiceClient
}

func NewUsersGrpc(client users.UsersServiceClient) UsersGrpc {
	return UsersGrpc{client: client}
}

func (s UsersGrpc) UpdateTrainingBalance(ctx context.Context, userID string, amountChange int) error {
	_, err := s.client.UpdateTrainingBalance(ctx, &users.UpdateTrainingBalanceRequest{
		UserId:       userID,
		AmountChange: int64(amountChange),
	})

	return err
}
