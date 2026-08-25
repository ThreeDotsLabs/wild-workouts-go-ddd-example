package training

import (
	"errors"
	"time"
)

// FreeCancellationPeriod is how long before the training it can still be canceled or
// rescheduled without losing credits.
const FreeCancellationPeriod = time.Hour * 24

func (t Training) CanBeCanceledForFree() bool {
	return time.Until(t.time) >= FreeCancellationPeriod
}

var ErrTrainingAlreadyCanceled = errors.New("training is already canceled")

func (t *Training) Cancel() error {
	if t.IsCanceled() {
		return ErrTrainingAlreadyCanceled
	}

	t.canceled = true
	return nil
}

func (t Training) IsCanceled() bool {
	return t.canceled
}
