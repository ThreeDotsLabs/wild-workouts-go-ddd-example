package training

import (
	"time"

	commonerrors "github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/errors"
	"github.com/pkg/errors"
)

type Training struct {
	uuid string

	userUUID string
	userName string

	time  time.Time
	notes string

	proposedNewTime time.Time
	moveProposedBy  UserType

	canceled bool
}

func NewTraining(uuid string, userUUID string, userName string, trainingTime time.Time) (*Training, error) {
	if uuid == "" {
		return nil, errors.New("empty training uuid")
	}
	if userUUID == "" {
		return nil, errors.New("empty userUUID")
	}
	if userName == "" {
		return nil, errors.New("empty userName")
	}
	if trainingTime.IsZero() {
		return nil, errors.New("zero training time")
	}

	return &Training{
		uuid:     uuid,
		userUUID: userUUID,
		userName: userName,
		time:     trainingTime,
	}, nil
}

// UnmarshalTrainingFromDatabase unmarshals Training from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalTrainingFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalTrainingFromDatabase(
	uuid string,
	userUUID string,
	userName string,
	trainingTime time.Time,
	notes string,
	canceled bool,
	proposedNewTime time.Time,
	moveProposedBy UserType,
) (*Training, error) {
	tr, err := NewTraining(uuid, userUUID, userName, trainingTime)
	if err != nil {
		return nil, err
	}

	tr.notes = notes
	tr.proposedNewTime = proposedNewTime
	tr.moveProposedBy = moveProposedBy
	tr.canceled = canceled

	return tr, nil
}

func (t Training) UUID() string {
	return t.uuid
}

func (t Training) UserUUID() string {
	return t.userUUID
}

func (t Training) UserName() string {
	return t.userName
}

func (t Training) Time() time.Time {
	return t.time
}

var ErrNoteTooLong = commonerrors.NewIncorrectInputError("Note too long", "note-too-long")

func (t *Training) UpdateNotes(notes string) error {
	if len(notes) > 1000 {
		return errors.WithStack(ErrNoteTooLong)
	}

	t.notes = notes
	return nil
}

func (t Training) Notes() string {
	return t.notes
}
