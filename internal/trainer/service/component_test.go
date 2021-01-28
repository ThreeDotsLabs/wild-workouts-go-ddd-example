package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	trainerHTTP "github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/server"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/tests"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/ports"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestHoursAvailability(t *testing.T) {
	t.Parallel()

	token := tests.FakeTrainerJWT(t, uuid.New().String())
	client := tests.NewTrainerHTTPClient(t, token)

	hour := tests.RelativeDate(11, 12)
	expectedHour := trainerHTTP.Hour{
		Available:            true,
		HasTrainingScheduled: false,
		Hour:                 hour,
	}

	date := hour.Truncate(24 * time.Hour)
	from := date.AddDate(0, 0, -1)
	to := date.AddDate(0, 0, 1)

	getHours := func() []trainerHTTP.Hour {
		dates := client.GetTrainerAvailableHours(t, from, to)
		for _, d := range dates {
			if d.Date.Equal(date) {
				return d.Hours
			}
		}
		t.Fatalf("Date not found in dates: %+v", dates)
		return nil
	}

	client.MakeHourUnavailable(t, hour)
	require.NotContains(t, getHours(), expectedHour)

	code := client.MakeHourAvailable(t, hour)
	require.Equal(t, http.StatusNoContent, code)
	require.Contains(t, getHours(), expectedHour)

	client.MakeHourUnavailable(t, hour)
	require.NotContains(t, getHours(), expectedHour)
}

func TestUnauthorizedForAttendee(t *testing.T) {
	t.Parallel()

	token := tests.FakeAttendeeJWT(t, uuid.New().String())
	client := tests.NewTrainerHTTPClient(t, token)

	hour := tests.RelativeDate(11, 13)

	code := client.MakeHourAvailable(t, hour)
	require.Equal(t, http.StatusUnauthorized, code)
}

func startService() bool {
	app := NewApplication(context.Background())

	trainerHTTPAddr := os.Getenv("TRAINER_HTTP_ADDR")
	go server.RunHTTPServerOnAddr(trainerHTTPAddr, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

	trainerGrpcAddr := os.Getenv("TRAINER_GRPC_ADDR")
	go server.RunGRPCServerOnAddr(trainerGrpcAddr, func(server *grpc.Server) {
		svc := ports.NewGrpcServer(app)
		trainer.RegisterTrainerServiceServer(server, svc)
	})

	ok := tests.WaitForPort(trainerHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for trainer HTTP to come up")
		return false
	}

	ok = tests.WaitForPort(trainerGrpcAddr)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		log.Println("Timed out waiting for trainings HTTP to come up")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
