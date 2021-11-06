package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/tests"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/ports"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateTraining(t *testing.T) {
	t.Parallel()

	token := tests.FakeAttendeeJWT(t, uuid.New().String())
	client := tests.NewTrainingsHTTPClient(t, token)

	hour := tests.RelativeDate(10, 12)
	trainingUUID := client.CreateTraining(t, "some note", hour)

	trainingsResponse := client.GetTrainings(t)

	var trainingsUUIDs []string
	for _, t := range trainingsResponse.Trainings {
		trainingsUUIDs = append(trainingsUUIDs, t.Uuid)
	}

	require.Contains(t, trainingsUUIDs, trainingUUID)
}

func TestCancelTraining(t *testing.T) {
	t.Parallel()

	token := tests.FakeAttendeeJWT(t, uuid.New().String())
	client := tests.NewTrainingsHTTPClient(t, token)

	hour := tests.RelativeDate(10, 13)
	trainingUUID := client.CreateTraining(t, "some note", hour)

	client.CancelTraining(t, trainingUUID, http.StatusOK)

	trainingsResponse := client.GetTrainings(t)

	var trainingsUUIDs []string
	for _, t := range trainingsResponse.Trainings {
		trainingsUUIDs = append(trainingsUUIDs, t.Uuid)
	}

	require.NotContains(t, trainingsUUIDs, trainingUUID)
}

func startService() bool {
	app := NewComponentTestApplication(context.Background())

	trainingsHTTPAddr := os.Getenv("TRAININGS_HTTP_ADDR")
	go server.RunHTTPServerOnAddr(trainingsHTTPAddr, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

	ok := tests.WaitForPort(trainingsHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for trainings HTTP to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
