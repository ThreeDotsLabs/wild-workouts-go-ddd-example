package training

import (
	"errors"
	"time"
)

func (t Training) CanBeCanceledForFree() bool {
	return time.Until(t.time) >= time.Hour*24
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
