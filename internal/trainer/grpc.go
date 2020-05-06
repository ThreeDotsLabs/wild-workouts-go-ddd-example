package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	db db
}

func (g GrpcServer) UpdateHour(ctx context.Context, req *trainer.UpdateHourRequest) (*trainer.EmptyResponse, error) {
	trainingTime, err := grpcTimestampToTime(req.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	date, err := g.db.DateModel(ctx, trainingTime)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to get data model: %s", err))
	}

	hour, found := date.FindHourInDate(trainingTime)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("%s hour not found in schedule", trainingTime))
	}

	if req.HasTrainingScheduled && !hour.Available {
		return nil, status.Error(codes.FailedPrecondition, "hour is not available for training")
	}

	if req.Available && req.HasTrainingScheduled {
		return nil, status.Error(codes.FailedPrecondition, "cannot set hour as available when it have training scheduled")
	}
	if !req.Available && !req.HasTrainingScheduled {
		return nil, status.Error(codes.FailedPrecondition, "cannot set hour as unavailable when it have no training scheduled")
	}
	hour.Available = req.Available

	if hour.HasTrainingScheduled && hour.HasTrainingScheduled == req.HasTrainingScheduled {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("hour HasTrainingScheduled is already %t", hour.HasTrainingScheduled))
	}

	hour.HasTrainingScheduled = req.HasTrainingScheduled
	if err := g.db.SaveModel(ctx, date); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to save date: %s", err))
	}

	return &trainer.EmptyResponse{}, nil
}

func (g GrpcServer) IsHourAvailable(ctx context.Context, req *trainer.IsHourAvailableRequest) (*trainer.IsHourAvailableResponse, error) {
	timeToCheck, err := grpcTimestampToTime(req.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	model, err := g.db.DateModel(ctx, timeToCheck)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("unable to get data model: %s", err))
	}

	if hour, found := model.FindHourInDate(timeToCheck); found {
		return &trainer.IsHourAvailableResponse{IsAvailable: hour.Available && !hour.HasTrainingScheduled}, nil
	}

	return &trainer.IsHourAvailableResponse{IsAvailable: false}, nil
}

func grpcTimestampToTime(timestamp *timestamp.Timestamp) (time.Time, error) {
	t, err := ptypes.Timestamp(timestamp)
	if err != nil {
		return time.Time{}, errors.New("unable to parse time")
	}

	t = t.UTC().Truncate(time.Hour)

	return t, nil
}
