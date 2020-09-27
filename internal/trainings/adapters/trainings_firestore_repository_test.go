package adapters_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/adapters"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app/query"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo - make tests parallel after fix of emulator: https://github.com/firebase/firebase-tools/issues/2452

func TestTrainingsFirestoreRepository_AddTraining(t *testing.T) {
	repo := newFirebaseRepository(t)

	testCases := []struct {
		Name                string
		TrainingConstructor func(t *testing.T) *training.Training
	}{
		{
			Name:                "standard_training",
			TrainingConstructor: newExampleTraining,
		},
		{
			Name:                "cancelled_training",
			TrainingConstructor: newCanceledTraining,
		},
		{
			Name:                "training_with_note",
			TrainingConstructor: newTrainingWithNote,
		},
		{
			Name:                "training_with_proposed_reschedule",
			TrainingConstructor: newTrainingWithProposedReschedule,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			ctx := context.Background()

			expectedTraining := c.TrainingConstructor(t)

			err := repo.AddTraining(ctx, expectedTraining)
			require.NoError(t, err)

			assertPersistedTrainingEquals(t, repo, expectedTraining)
		})
	}
}

func TestTrainingsFirestoreRepository_UpdateTraining(t *testing.T) {
	repo := newFirebaseRepository(t)
	ctx := context.Background()

	expectedTraining := newExampleTraining(t)

	err := repo.AddTraining(ctx, expectedTraining)
	require.NoError(t, err)

	var updatedTraining *training.Training

	err = repo.UpdateTraining(
		ctx,
		expectedTraining.UUID(),
		training.MustNewUser(expectedTraining.UserUUID(), training.Attendee),
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			assertTrainingsEquals(t, expectedTraining, tr)

			err := tr.UpdateNotes("note")
			require.NoError(t, err)

			updatedTraining = tr

			return tr, nil
		},
	)
	require.NoError(t, err)

	assertPersistedTrainingEquals(t, repo, updatedTraining)
}

func TestTrainingsFirestoreRepository_GetTraining_not_exists(t *testing.T) {
	repo := newFirebaseRepository(t)

	trainingUUID := uuid.New().String()

	tr, err := repo.GetTraining(
		context.Background(),
		trainingUUID,
		training.MustNewUser(uuid.New().String(), training.Attendee),
	)
	assert.Nil(t, tr)
	assert.EqualError(t, err, training.NotFoundError{trainingUUID}.Error())
}

func TestTrainingsFirestoreRepository_get_and_update_another_users_training(t *testing.T) {
	repo := newFirebaseRepository(t)

	ctx := context.Background()
	tr := newExampleTraining(t)

	err := repo.AddTraining(ctx, tr)
	require.NoError(t, err)

	assertPersistedTrainingEquals(t, repo, tr)

	requestingUser := training.MustNewUser(uuid.New().String(), training.Attendee)

	_, err = repo.GetTraining(
		context.Background(),
		tr.UUID(),
		requestingUser,
	)
	assert.EqualError(
		t,
		err,
		training.ForbiddenToSeeTrainingError{
			RequestingUserUUID: requestingUser.UUID(),
			TrainingOwnerUUID:  tr.UserUUID(),
		}.Error(),
	)

	err = repo.UpdateTraining(
		ctx,
		tr.UUID(),
		requestingUser,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			return nil, nil
		},
	)
	assert.EqualError(
		t,
		err,
		training.ForbiddenToSeeTrainingError{
			RequestingUserUUID: requestingUser.UUID(),
			TrainingOwnerUUID:  tr.UserUUID(),
		}.Error(),
	)
}

func TestTrainingsFirestoreRepository_AllTrainings(t *testing.T) {
	repo := newFirebaseRepository(t)

	// AllTrainings returns all documents, because of that we need to do exception and do DB cleanup
	// In general, I recommend to do it before test. In that way you are sure that cleanup is done.
	// Thanks to that tests are more stable.
	// More about why it is important you can find in https://threedots.tech/post/database-integration-testing/
	err := repo.RemoveAllTrainings(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	exampleTraining := newExampleTraining(t)
	canceledTraining := newCanceledTraining(t)
	trainingWithNote := newTrainingWithNote(t)
	trainingWithProposedReschedule := newTrainingWithProposedReschedule(t)

	trainingsToAdd := []*training.Training{
		exampleTraining,
		canceledTraining,
		trainingWithNote,
		trainingWithProposedReschedule,
	}

	for _, tr := range trainingsToAdd {
		err = repo.AddTraining(ctx, tr)
		require.NoError(t, err)
	}

	trainings, err := repo.AllTrainings(context.Background())
	require.NoError(t, err)

	proposedNewTime := trainingWithProposedReschedule.ProposedNewTime()
	proposer := trainingWithProposedReschedule.MovedProposedBy().String()

	expectedTrainings := []query.Training{
		{
			UUID:           exampleTraining.UUID(),
			UserUUID:       exampleTraining.UserUUID(),
			User:           "User",
			Time:           exampleTraining.Time(),
			Notes:          "",
			CanBeCancelled: true,
		},
		{
			UUID:           trainingWithNote.UUID(),
			UserUUID:       trainingWithNote.UserUUID(),
			User:           "User",
			Time:           trainingWithNote.Time(),
			Notes:          trainingWithNote.Notes(),
			CanBeCancelled: true,
		},
		{
			UUID:           trainingWithProposedReschedule.UUID(),
			UserUUID:       trainingWithProposedReschedule.UserUUID(),
			User:           "User",
			Time:           trainingWithProposedReschedule.Time(),
			Notes:          "",
			ProposedTime:   &proposedNewTime,
			MoveProposedBy: &proposer,
			CanBeCancelled: true,
		},
	}

	assertQueryTrainingsEquals(t, expectedTrainings, trainings)
}

func TestTrainingsFirestoreRepository_FindTrainingsForUser(t *testing.T) {
	repo := newFirebaseRepository(t)

	ctx := context.Background()

	userUUID := uuid.New().String()

	tr1, err := training.NewTraining(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	err = repo.AddTraining(ctx, tr1)
	require.NoError(t, err)

	tr2, err := training.NewTraining(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	err = repo.AddTraining(ctx, tr2)
	require.NoError(t, err)

	// this training should be not in the list
	canceledTraining, err := training.NewTraining(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	err = canceledTraining.Cancel()
	require.NoError(t, err)
	err = repo.AddTraining(ctx, canceledTraining)
	require.NoError(t, err)

	trainings, err := repo.FindTrainingsForUser(context.Background(), userUUID)
	require.NoError(t, err)

	assertQueryTrainingsEquals(t, trainings, []query.Training{
		{
			UUID:     tr1.UUID(),
			UserUUID: userUUID,
			User:     "User",
			Time:     tr1.Time(),
		},
		{
			UUID:     tr2.UUID(),
			UserUUID: userUUID,
			User:     "User",
			Time:     tr2.Time(),
		},
	})
}

func newRandomTrainingTime() time.Time {
	min := time.Now().AddDate(0, 0, 5).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func newExampleTraining(t *testing.T) *training.Training {
	tr, err := training.NewTraining(
		uuid.New().String(),
		uuid.New().String(),
		"User",
		newRandomTrainingTime(),
	)
	require.NoError(t, err)

	return tr
}

func newCanceledTraining(t *testing.T) *training.Training {
	tr, err := training.NewTraining(
		uuid.New().String(),
		uuid.New().String(),
		"User",
		newRandomTrainingTime(),
	)
	require.NoError(t, err)

	err = tr.Cancel()
	require.NoError(t, err)

	return tr
}

func newTrainingWithNote(t *testing.T) *training.Training {
	tr := newExampleTraining(t)
	err := tr.UpdateNotes("foo")
	require.NoError(t, err)

	return tr
}

func newTrainingWithProposedReschedule(t *testing.T) *training.Training {
	tr := newExampleTraining(t)
	tr.ProposeReschedule(time.Now().AddDate(0, 0, 14), training.Trainer)

	return tr
}

func assertPersistedTrainingEquals(t *testing.T, repo adapters.TrainingsFirestoreRepository, tr *training.Training) {
	persistedTraining, err := repo.GetTraining(
		context.Background(),
		tr.UUID(),
		training.MustNewUser(tr.UserUUID(), training.Attendee),
	)
	require.NoError(t, err)

	assertTrainingsEquals(t, tr, persistedTraining)
}

// Firestore is not storing time with same precision, so we need to round it a bit
var cmpRoundTimeOpt = cmp.Comparer(func(x, y time.Time) bool {
	return x.Truncate(time.Millisecond).Equal(y.Truncate(time.Millisecond))
})

func assertTrainingsEquals(t *testing.T, tr1, tr2 *training.Training) {
	cmpOpts := []cmp.Option{
		cmpRoundTimeOpt,
		cmp.AllowUnexported(
			training.UserType{},
			time.Time{},
			training.Training{},
		),
	}

	assert.True(
		t,
		cmp.Equal(tr1, tr2, cmpOpts...),
		cmp.Diff(tr1, tr2, cmpOpts...),
	)
}

func assertQueryTrainingsEquals(t *testing.T, expectedTrainings, trainings []query.Training) bool {
	cmpOpts := []cmp.Option{
		cmpRoundTimeOpt,
		cmpopts.SortSlices(func(x, y query.Training) bool {
			return x.Time.After(y.Time)
		}),
	}
	return assert.True(t,
		cmp.Equal(expectedTrainings, trainings, cmpOpts...),
		cmp.Diff(expectedTrainings, trainings, cmpOpts...),
	)
}

func newFirebaseRepository(t *testing.T) adapters.TrainingsFirestoreRepository {
	firestoreClient, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	require.NoError(t, err)

	return adapters.NewTrainingsFirestoreRepository(firestoreClient)
}
