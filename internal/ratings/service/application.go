package service

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/ratings/app"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	return newApplication(ctx), func() {}
}

func NewComponentTestApplication(ctx context.Context) app.Application {
	return newApplication(ctx)
}

func newApplication(ctx context.Context) app.Application {
	return app.Application{
		// TODO: inject commands
		// you can check internal/trainings/service/application.go for the inspiration
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
