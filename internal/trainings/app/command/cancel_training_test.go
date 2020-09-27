package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app/command"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCancelTraining(t *testing.T) {
	requestingUserID := "requesting-user-id"

	testCases := []struct {
		Name     string
		UserType training.UserType

		TrainingConstructor func() *training.Training

		ShouldFail    bool
		ExpectedError string

		ShouldUpdateBalance   bool
		ExpectedBalanceChange int
	}{
		{
			Name:     "return_training_balance_when_attendee_cancels",
			UserType: training.Attendee,
			TrainingConstructor: func() *training.Training {
				return createExampleTraining(t, requestingUserID, time.Now().Add(48*time.Hour))
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 1,
		},
		{
			Name:     "return_training_balance_when_trainer_cancels",
			UserType: training.Trainer,
			TrainingConstructor: func() *training.Training {
				return createExampleTraining(t, "trainer-id", time.Now().Add(48*time.Hour))
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 1,
		},
		{
			Name:     "extra_training_balance_when_trainer_cancels_before_24h",
			UserType: training.Trainer,
			TrainingConstructor: func() *training.Training {
				return createExampleTraining(t, "trainer-id", time.Now().Add(12*time.Hour))
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 2,
		},
		{
			Name:     "no_training_balance_returned_when_attendee_cancels_before_24h",
			UserType: training.Attendee,
			TrainingConstructor: func() *training.Training {
				return createExampleTraining(t, requestingUserID, time.Now().Add(12*time.Hour))
			},
			ShouldUpdateBalance: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			trainingUUID := "any-training-uuid"
			deps := newDependencies()

			tr := tc.TrainingConstructor()
			deps.repository.Trainings = map[string]training.Training{
				trainingUUID: *tr,
			}

			err := deps.handler.Handle(context.Background(), command.CancelTraining{
				TrainingUUID: trainingUUID,
				User:         training.MustNewUser(requestingUserID, tc.UserType),
			})

			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)

			if tc.ShouldUpdateBalance {
				require.Len(t, deps.userService.balanceUpdates, 1)
				require.Equal(t, tr.UserUUID(), deps.userService.balanceUpdates[0].userID)
				require.Equal(t, tc.ExpectedBalanceChange, deps.userService.balanceUpdates[0].amountChange)
			} else {
				require.Len(t, deps.userService.balanceUpdates, 0)
			}

			require.Len(t, deps.trainerService.trainingsCancelled, 1)
			require.Equal(t, tr.Time(), deps.trainerService.trainingsCancelled[0])
		})
	}
}

func createExampleTraining(t *testing.T, requestingUserID string, trainingTime time.Time) *training.Training {
	tr, err := training.NewTraining(
		uuid.New().String(),
		requestingUserID,
		"foo",
		trainingTime,
	)
	require.NoError(t, err)

	return tr
}

type dependencies struct {
	repository     *repositoryMock
	trainerService *trainerServiceMock
	userService    *userServiceMock
	handler        command.CancelTrainingHandler
}

func newDependencies() dependencies {
	repository := &repositoryMock{}
	trainerService := &trainerServiceMock{}
	userService := &userServiceMock{}

	return dependencies{
		repository:     repository,
		trainerService: trainerService,
		userService:    userService,
		handler:        command.NewCancelTrainingHandler(repository, userService, trainerService),
	}
}

type repositoryMock struct {
	Trainings map[string]training.Training
}

func (r *repositoryMock) GetTraining(ctx context.Context, trainingUUID string, user training.User) (*training.Training, error) {
	panic("implement me")
}

func (r *repositoryMock) UpdateTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
	updateFn func(ctx context.Context, tr *training.Training) (*training.Training, error),
) error {
	tr, ok := r.Trainings[trainingUUID]
	if !ok {
		return errors.Errorf("training '%s' not found", trainingUUID)
	}

	updatedTraining, err := updateFn(ctx, &tr)
	if err != nil {
		return err
	}

	r.Trainings[trainingUUID] = *updatedTraining

	return nil
}

func (r repositoryMock) AddTraining(ctx context.Context, tr *training.Training) error {
	panic("implement me")
}

type trainerServiceMock struct {
	trainingsCancelled []time.Time
}

func (t *trainerServiceMock) MoveTraining(ctx context.Context, newTime time.Time, originalTrainingTime time.Time) error {
	panic("implement me")
}

func (t *trainerServiceMock) ScheduleTraining(ctx context.Context, trainingTime time.Time) error {
	panic("implement me")
}

func (t *trainerServiceMock) CancelTraining(ctx context.Context, trainingTime time.Time) error {
	t.trainingsCancelled = append(t.trainingsCancelled, trainingTime)
	return nil
}

type balanceUpdate struct {
	userID       string
	amountChange int
}

type userServiceMock struct {
	balanceUpdates []balanceUpdate
}

func (u *userServiceMock) UpdateTrainingBalance(ctx context.Context, userID string, amountChange int) error {
	u.balanceUpdates = append(u.balanceUpdates, balanceUpdate{userID, amountChange})
	return nil
}
