package training_test

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraining_RescheduleTraining(t *testing.T) {
	t.Parallel()
	tr := newExampleTraining(t)

	oldTime := tr.Time()
	newTime := time.Now().AddDate(0, 0, 14).Round(time.Hour)

	// it's always a good idea to ensure about pre-conditions in the test ;-)
	assert.False(t, oldTime.Equal(newTime))

	err := tr.RescheduleTraining(newTime)
	assert.NoError(t, err)
	assert.True(t, tr.Time().Equal(newTime))
}

func TestTraining_RescheduleTraining_less_than_24h_before(t *testing.T) {
	t.Parallel()
	originalTime := time.Now().Round(time.Hour)
	rescheduleRequestTime := originalTime.AddDate(0, 0, 5)

	tr := newExampleTrainingWithTime(t, originalTime)

	err := tr.RescheduleTraining(rescheduleRequestTime)

	assert.EqualError(t, err, training.CantRescheduleBeforeTimeError{
		TrainingTime: tr.Time(),
	}.Error())
}

func TestTraining_ProposeReschedule_by_attendee(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Proposer training.UserType
		Approver training.UserType
	}{
		{
			Name:     "proposed_by_attendee",
			Proposer: training.Attendee,
			Approver: training.Trainer,
		},
		{
			Name:     "proposed_by_trainer",
			Proposer: training.Trainer,
			Approver: training.Attendee,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			originalTime := time.Now().Round(time.Hour)
			rescheduleRequestTime := originalTime.AddDate(0, 0, 5)
			tr := newExampleTrainingWithTime(t, originalTime)

			assert.False(t, tr.IsRescheduleProposed())

			tr.ProposeReschedule(rescheduleRequestTime, c.Proposer)

			assert.True(t, tr.IsRescheduleProposed())

			err := tr.ApproveReschedule(c.Approver)
			require.NoError(t, err)

			tr.Time().Equal(rescheduleRequestTime)
			assert.False(t, tr.IsRescheduleProposed())
		})
	}
}

func TestTraining_ProposeReschedule_approve_by_proposer(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Proposer training.UserType
	}{
		{
			Proposer: training.Attendee,
		},
		{
			Proposer: training.Trainer,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Proposer.String(), func(t *testing.T) {
			t.Parallel()
			originalTime := time.Now().Round(time.Hour)
			rescheduleRequestTime := originalTime.AddDate(0, 0, 5)
			tr := newExampleTrainingWithTime(t, originalTime)

			tr.ProposeReschedule(rescheduleRequestTime, c.Proposer)

			err := tr.ApproveReschedule(c.Proposer)
			assert.Error(t, err)

			tr.Time().Equal(originalTime)
			assert.True(t, tr.IsRescheduleProposed())
		})
	}
}

func TestTraining_ApproveReschedule_not_proposed(t *testing.T) {
	t.Parallel()
	tr := newExampleTrainingWithTime(t, time.Now().Round(time.Hour))

	assert.EqualError(t, tr.ApproveReschedule(training.Trainer), training.ErrNoRescheduleRequested.Error())
}

func TestTraining_RejectRescheduleTraining(t *testing.T) {
	t.Parallel()
	originalTime := time.Now().Round(time.Hour)
	rescheduleRequestTime := originalTime.AddDate(0, 0, 5)
	tr := newExampleTrainingWithTime(t, originalTime)

	tr.ProposeReschedule(rescheduleRequestTime, training.Attendee)

	err := tr.RejectReschedule()
	assert.NoError(t, err)

	tr.Time().Equal(originalTime)
	assert.False(t, tr.IsRescheduleProposed())
}
