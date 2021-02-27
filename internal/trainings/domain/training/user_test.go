package training_test

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUserAllowedToSeeTraining(t *testing.T) {
	attendee1, err := training.NewUser(uuid.New().String(), training.Attendee)
	require.NoError(t, err)

	attendee2, err := training.NewUser(uuid.New().String(), training.Attendee)
	require.NoError(t, err)

	trainer, err := training.NewUser(uuid.New().String(), training.Trainer)
	require.NoError(t, err)

	testCases := []struct {
		Name              string
		CreateTraining    func(t *testing.T) *training.Training
		User              training.User
		ExpectedIsAllowed bool
	}{
		{
			Name: "attendees_training",
			CreateTraining: func(t *testing.T) *training.Training {
				tr, err := training.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              attendee1,
			ExpectedIsAllowed: true,
		},
		{
			Name: "another_attendees_training",
			CreateTraining: func(t *testing.T) *training.Training {
				tr, err := training.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              attendee2,
			ExpectedIsAllowed: false,
		},
		{
			Name: "trainer",
			CreateTraining: func(t *testing.T) *training.Training {
				tr, err := training.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              trainer,
			ExpectedIsAllowed: true, // trainer have access to all trainings
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tr := c.CreateTraining(t)

			err := training.CanUserSeeTraining(c.User, *tr)

			if c.ExpectedIsAllowed {

			} else {
				assert.EqualError(
					t,
					err,
					training.ForbiddenToSeeTrainingError{c.User.UUID(), tr.UserUUID()}.Error(),
				)
			}
		})
	}
}
