package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/tests"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/ratings/ports"
	"github.com/go-chi/chi/v5"
)

// TODO: implement component test - https://threedots.tech/post/microservices-test-architecture/

func startService() bool {
	app := NewComponentTestApplication(context.Background())

	ratingsHTTPAddr := os.Getenv("RATINGS_HTTP_ADDR")
	go server.RunHTTPServerOnAddr(ratingsHTTPAddr, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

	ok := tests.WaitForPort(ratingsHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for ratings HTTP to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
