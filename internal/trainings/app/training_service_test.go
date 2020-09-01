package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app"
	"github.com/stretchr/testify/require"
)

func TestCancelTraining(t *testing.T) {
	requestingUserID := "requesting-user-id"

	testCases := []struct {
		Name     string
		UserRole string

		Training app.Training

		ShouldFail    bool
		ExpectedError string

		ShouldUpdateBalance   bool
		ExpectedBalanceChange int
	}{
		{
			Name:     "return_training_balance_when_attendee_cancels",
			UserRole: "attendee",
			Training: app.Training{
				UserUUID: requestingUserID,
				Time:     time.Now().Add(48 * time.Hour),
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 1,
		},
		{
			Name:     "return_training_balance_when_trainer_cancels",
			UserRole: "trainer",
			Training: app.Training{
				UserUUID: "trainer-id",
				Time:     time.Now().Add(48 * time.Hour),
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 1,
		},
		{
			Name:     "extra_training_balance_when_trainer_cancels_before_24h",
			UserRole: "trainer",
			Training: app.Training{
				UserUUID: "trainer-id",
				Time:     time.Now().Add(12 * time.Hour),
			},
			ShouldUpdateBalance:   true,
			ExpectedBalanceChange: 2,
		},
		{
			Name:     "no_training_balance_returned_when_attendee_cancels_before_24h",
			UserRole: "attendee",
			Training: app.Training{
				UserUUID: requestingUserID,
				Time:     time.Now().Add(12 * time.Hour),
			},
			ShouldUpdateBalance: false,
		},
		{
			Name:     "fail_updating_other_attendee_training",
			UserRole: "attendee",
			Training: app.Training{
				UserUUID: "another-attendee-id",
				Time:     time.Now().Add(48 * time.Hour),
			},
			ShouldFail:          true,
			ExpectedError:       "user 'requesting-user-id' is trying to cancel training of user 'another-attendee-id'",
			ShouldUpdateBalance: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			trainingUUID := "any-training-uuid"
			deps := newDependencies()
			deps.repository.training = tc.Training

			user := auth.User{
				UUID: requestingUserID,
				Role: tc.UserRole,
			}

			err := deps.trainingsService.CancelTraining(context.Background(), user, trainingUUID)

			if tc.ShouldFail {
				require.EqualError(t, err, tc.ExpectedError)
				return
			}

			require.NoError(t, err)

			if tc.ShouldUpdateBalance {
				require.Len(t, deps.userService.balanceUpdates, 1)
				require.Equal(t, tc.Training.UserUUID, deps.userService.balanceUpdates[0].userID)
				require.Equal(t, tc.ExpectedBalanceChange, deps.userService.balanceUpdates[0].amountChange)
			} else {
				require.Len(t, deps.userService.balanceUpdates, 0)
			}

			require.Len(t, deps.trainerService.trainingsCancelled, 1)
			require.Equal(t, tc.Training.Time, deps.trainerService.trainingsCancelled[0])
		})
	}
}

type dependencies struct {
	repository       *repositoryMock
	trainerService   *trainerServiceMock
	userService      *userServiceMock
	trainingsService app.TrainingService
}

func newDependencies() dependencies {
	repository := &repositoryMock{}
	trainerService := &trainerServiceMock{}
	userService := &userServiceMock{}

	return dependencies{
		repository:       repository,
		trainerService:   trainerService,
		userService:      userService,
		trainingsService: app.NewTrainingsService(repository, trainerService, userService),
	}
}

type repositoryMock struct {
	training app.Training
}

func (r repositoryMock) FindTrainingsForUser(ctx context.Context, user auth.User) ([]app.Training, error) {
	panic("implement me")
}

func (r repositoryMock) AllTrainings(ctx context.Context) ([]app.Training, error) {
	panic("implement me")
}

func (r repositoryMock) CreateTraining(ctx context.Context, training app.Training, createFn func() error) error {
	panic("implement me")
}

func (r repositoryMock) CancelTraining(ctx context.Context, trainingUUID string, deleteFn func(app.Training) error) error {
	return deleteFn(r.training)
}

func (r repositoryMock) RescheduleTraining(ctx context.Context, trainingUUID string, newTime time.Time, updateFn func(app.Training) (app.Training, error)) error {
	panic("implement me")
}

func (r repositoryMock) ApproveTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	panic("implement me")
}

func (r repositoryMock) RejectTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	panic("implement me")
}

type trainerServiceMock struct {
	trainingsCancelled []time.Time
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
