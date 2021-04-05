package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client/trainings"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client/users"
	"github.com/stretchr/testify/require"
)

func authorizationBearer(token string) func(context.Context, *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
}

type TrainerHTTPClient struct {
	client *trainer.ClientWithResponses
}

func NewTrainerHTTPClient(t *testing.T, token string) TrainerHTTPClient {
	addr := os.Getenv("TRAINER_HTTP_ADDR")
	ok := WaitForPort(addr)
	require.True(t, ok, "Trainer HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := trainer.NewClientWithResponses(
		url,
		trainer.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return TrainerHTTPClient{
		client: client,
	}
}

func (c TrainerHTTPClient) MakeHourAvailable(t *testing.T, hour time.Time) int {
	response, err := c.client.MakeHourAvailable(context.Background(), trainer.MakeHourAvailableJSONRequestBody{
		Hours: []time.Time{hour},
	})
	require.NoError(t, err)
	return response.StatusCode
}

func (c TrainerHTTPClient) MakeHourUnavailable(t *testing.T, hour time.Time) {
	response, err := c.client.MakeHourUnavailable(context.Background(), trainer.MakeHourUnavailableJSONRequestBody{
		Hours: []time.Time{hour},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
}

func (c TrainerHTTPClient) GetTrainerAvailableHours(t *testing.T, from time.Time, to time.Time) []trainer.Date {
	response, err := c.client.GetTrainerAvailableHoursWithResponse(context.Background(), &trainer.GetTrainerAvailableHoursParams{
		DateFrom: from,
		DateTo:   to,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

type TrainingsHTTPClient struct {
	client *trainings.ClientWithResponses
}

func NewTrainingsHTTPClient(t *testing.T, token string) TrainingsHTTPClient {
	addr := os.Getenv("TRAININGS_HTTP_ADDR")
	fmt.Println("Trying trainings http:", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Trainings HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := trainings.NewClientWithResponses(
		url,
		trainings.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return TrainingsHTTPClient{
		client: client,
	}
}

func (c TrainingsHTTPClient) CreateTraining(t *testing.T, note string, hour time.Time) string {
	response, err := c.client.CreateTrainingWithResponse(context.Background(), trainings.CreateTrainingJSONRequestBody{
		Notes: note,
		Time:  hour,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

func (c TrainingsHTTPClient) CreateTrainingShouldFail(t *testing.T, note string, hour time.Time) {
	response, err := c.client.CreateTraining(context.Background(), trainings.CreateTrainingJSONRequestBody{
		Notes: note,
		Time:  hour,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func (c TrainingsHTTPClient) GetTrainings(t *testing.T) trainings.Trainings {
	response, err := c.client.GetTrainingsWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c TrainingsHTTPClient) CancelTraining(t *testing.T, trainingUUID string, expectedStatusCode int) {
	response, err := c.client.CancelTraining(context.Background(), trainingUUID)
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, response.StatusCode)
}

type UsersHTTPClient struct {
	client *users.ClientWithResponses
}

func NewUsersHTTPClient(t *testing.T, token string) UsersHTTPClient {
	addr := os.Getenv("USERS_HTTP_ADDR")
	ok := WaitForPort(addr)
	require.True(t, ok, "Users HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := users.NewClientWithResponses(
		url,
		users.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return UsersHTTPClient{
		client: client,
	}
}

func (c UsersHTTPClient) GetCurrentUser(t *testing.T) users.User {
	response, err := c.client.GetCurrentUserWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func lastPathElement(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
