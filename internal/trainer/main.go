package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	firebaseDB := db{firestoreClient}

	hourFactory, err := hour.NewFactory(hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	})
	if err != nil {
		panic(err)
	}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(firebaseDB)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return HandlerFromMux(
				HttpServer{
					firebaseDB,
					NewFirestoreHourRepository(firestoreClient, hourFactory),
				},
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := GrpcServer{NewFirestoreHourRepository(firestoreClient, hourFactory)}
			trainer.RegisterTrainerServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
