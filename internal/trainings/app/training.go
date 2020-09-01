package app

import "time"

type Training struct {
	UUID     string
	UserUUID string
	User     string

	Time  time.Time
	Notes string

	ProposedTime   *time.Time
	MoveProposedBy *string
}

func (t Training) CanBeCancelled() bool {
	return t.Time.Sub(time.Now()) > time.Hour*24
}

func (t Training) MoveRequiresAccept() bool {
	return !t.CanBeCancelled()
}
