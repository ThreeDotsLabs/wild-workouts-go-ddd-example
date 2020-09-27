package query

import "time"

type Training struct {
	UUID     string
	UserUUID string
	User     string

	Time  time.Time
	Notes string

	ProposedTime   *time.Time
	MoveProposedBy *string

	CanBeCancelled bool
}
