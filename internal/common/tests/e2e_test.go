package tests

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateTraining(t *testing.T) {
	t.Parallel()

	hour := RelativeDate(12, 12)

	userID := "TestCreateTraining-user"
	trainerJWT := FakeTrainerJWT(t, uuid.New().String())
	attendeeJWT := FakeAttendeeJWT(t, userID)
	trainerHTTPClient := NewTrainerHTTPClient(t, trainerJWT)
	trainingsHTTPClient := NewTrainingsHTTPClient(t, attendeeJWT)
	usersHTTPClient := NewUsersHTTPClient(t, attendeeJWT)

	usersGrpcClient, _, err := client.NewUsersClient()
	require.NoError(t, err)

	// Cancel the training if exists and make the hour available
	trainings := trainingsHTTPClient.GetTrainings(t)
	for _, training := range trainings.Trainings {
		if training.Time.Equal(hour) {
			trainingsTrainerHTTPClient := NewTrainingsHTTPClient(t, trainerJWT)
			trainingsTrainerHTTPClient.CancelTraining(t, training.Uuid, 200)
			break
		}
	}
	hours := trainerHTTPClient.GetTrainerAvailableHours(t, hour, hour)
	if len(hours) > 0 {
		for _, h := range hours[0].Hours {
			if h.Hour.Equal(hour) {
				trainerHTTPClient.MakeHourUnavailable(t, hour)
				break
			}
		}
	}

	trainerHTTPClient.MakeHourAvailable(t, hour)

	user := usersHTTPClient.GetCurrentUser(t)
	originalBalance := user.Balance

	_, err = usersGrpcClient.UpdateTrainingBalance(context.Background(), &users.UpdateTrainingBalanceRequest{
		UserId:       userID,
		AmountChange: 1,
	})
	require.NoError(t, err)

	user = usersHTTPClient.GetCurrentUser(t)
	require.Equal(t, originalBalance+1, user.Balance, "Attendee's balance should be updated")

	trainingUUID := trainingsHTTPClient.CreateTraining(t, "some note", hour)

	trainingsResponse := trainingsHTTPClient.GetTrainings(t)
	require.Len(t, trainingsResponse.Trainings, 1)
	require.Equal(t, trainingUUID, trainingsResponse.Trainings[0].Uuid, "Attendee should see the training")

	user = usersHTTPClient.GetCurrentUser(t)
	require.Equal(t, originalBalance, user.Balance, "Attendee's balance should be updated after a training is scheduled")
}
