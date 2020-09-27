package training_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTraining(t *testing.T) {
	trainingUUID := uuid.New().String()
	userUUID := uuid.New().String()
	userName := "user name"
	trainingTime := time.Now().Round(time.Hour)

	tr, err := training.NewTraining(trainingUUID, userUUID, userName, trainingTime)
	require.NoError(t, err)

	assert.Equal(t, trainingUUID, tr.UUID())
	assert.Equal(t, userUUID, tr.UserUUID())
	assert.Equal(t, trainingTime, tr.Time())
	assert.Equal(t, userName, tr.UserName())
}

func TestNewTraining_invalid(t *testing.T) {
	trainingUUID := uuid.New().String()
	userUUID := uuid.New().String()
	trainingTime := time.Now().Round(time.Hour)
	userName := "user name"

	_, err := training.NewTraining("", userUUID, userName, trainingTime)
	assert.Error(t, err)

	_, err = training.NewTraining(trainingUUID, "", userName, trainingTime)
	assert.Error(t, err)

	_, err = training.NewTraining(trainingUUID, userUUID, userName, time.Time{})
	assert.Error(t, err)

	_, err = training.NewTraining(trainingUUID, userUUID, "", time.Time{})
	assert.Error(t, err)
}

func TestTraining_UpdateNotes(t *testing.T) {
	tr := newExampleTraining(t)
	// it's always a good idea to ensure about pre-conditions in the test ;-)
	require.Equal(t, "", tr.Notes())

	err := tr.UpdateNotes("foo")
	require.NoError(t, err)
	assert.Equal(t, "foo", tr.Notes())
}

func TestTraining_UpdateNotes_too_long(t *testing.T) {
	tr := newExampleTraining(t)

	err := tr.UpdateNotes(strings.Repeat("x", 1001))
	assert.EqualError(t, err, training.ErrNoteTooLong.Error())
}

func TestTraining_MoreThanDayUntilTraining(t *testing.T) {
	trainingNow := newExampleTrainingWithTime(t, time.Now())
	assert.False(t, trainingNow.CanBeCanceledForFree())

	trainingInTwoDays := newExampleTrainingWithTime(t, time.Now().AddDate(0, 0, 2))
	assert.True(t, trainingInTwoDays.CanBeCanceledForFree())
}

func newExampleTraining(t *testing.T) *training.Training {
	tr, err := training.NewTraining(
		uuid.New().String(),
		uuid.New().String(),
		"user name",
		time.Now().AddDate(0, 0, 5).Round(time.Hour),
	)
	require.NoError(t, err)

	return tr
}

func newExampleTrainingWithTime(t *testing.T, trainingTime time.Time) *training.Training {
	tr, err := training.NewTraining(
		uuid.New().String(),
		uuid.New().String(),
		"user name",
		trainingTime,
	)
	require.NoError(t, err)

	return tr
}

func newCanceledTraining(t *testing.T) *training.Training {
	tr := newExampleTraining(t)
	require.NoError(t, tr.Cancel())

	return tr
}
