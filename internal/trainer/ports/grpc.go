package ports

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) MakeHourAvailable(ctx context.Context, request *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime := protoTimestampToTime(request.Time)

	if err := g.app.Commands.MakeHoursAvailable.Handle(ctx, []time.Time{trainingTime}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) ScheduleTraining(ctx context.Context, request *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime := protoTimestampToTime(request.Time)

	if err := g.app.Commands.ScheduleTraining.Handle(ctx, trainingTime); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) CancelTraining(ctx context.Context, request *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime := protoTimestampToTime(request.Time)

	if err := g.app.Commands.CancelTraining.Handle(ctx, trainingTime); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) IsHourAvailable(ctx context.Context, request *trainer.IsHourAvailableRequest) (*trainer.IsHourAvailableResponse, error) {
	trainingTime := protoTimestampToTime(request.Time)

	isAvailable, err := g.app.Queries.HourAvailability.Handle(ctx, trainingTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &trainer.IsHourAvailableResponse{IsAvailable: isAvailable}, nil
}

func protoTimestampToTime(timestamp *timestamp.Timestamp) time.Time {
	return timestamp.AsTime().UTC().Truncate(time.Hour)
}
