package training

import (
	"fmt"
	"time"

	commonErrors "github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/errors"
	"github.com/pkg/errors"
)

func (t Training) MovedProposedBy() UserType {
	return t.moveProposedBy
}

func (t Training) ProposedNewTime() time.Time {
	return t.proposedNewTime
}

type CantRescheduleBeforeTimeError struct {
	TrainingTime time.Time
}

func (c CantRescheduleBeforeTimeError) Error() string {
	return fmt.Sprintf(
		"Training can be rescheduled up to %.0f hours before it starts (training time: %s)",
		FreeCancellationPeriod.Hours(),
		c.TrainingTime.Format("2006-01-02 15:04 MST"),
	)
}

func (c CantRescheduleBeforeTimeError) Slug() string {
	return "training-reschedule-too-late"
}

func (c CantRescheduleBeforeTimeError) ErrorType() commonErrors.ErrorType {
	return commonErrors.ErrorTypeIncorrectInput
}

func (t *Training) RescheduleTraining(newTime time.Time) error {
	if !t.CanBeCanceledForFree() {
		err := CantRescheduleBeforeTimeError{
			TrainingTime: t.Time(),
		}
		return errors.WithStack(err)
	}

	t.time = newTime

	return nil
}

func (t *Training) ProposeReschedule(newTime time.Time, proposerType UserType) {
	t.moveProposedBy = proposerType
	t.proposedNewTime = newTime
}

func (t *Training) IsRescheduleProposed() bool {
	return !t.moveProposedBy.IsZero() && !t.proposedNewTime.IsZero()
}

var ErrNoRescheduleRequested = errors.New("no training reschedule was requested yet")

func (t *Training) ApproveReschedule(userType UserType) error {
	if !t.IsRescheduleProposed() {
		return errors.WithStack(ErrNoRescheduleRequested)
	}

	if t.moveProposedBy == userType {
		return errors.Errorf(
			"trying to approve reschedule by the same user type which proposed reschedule (%s)",
			userType.String(),
		)
	}

	t.time = t.proposedNewTime

	t.proposedNewTime = time.Time{}
	t.moveProposedBy = UserType{}

	return nil
}

func (t *Training) RejectReschedule() error {
	if !t.IsRescheduleProposed() {
		return errors.WithStack(ErrNoRescheduleRequested)
	}

	t.proposedNewTime = time.Time{}
	t.moveProposedBy = UserType{}

	return nil
}
