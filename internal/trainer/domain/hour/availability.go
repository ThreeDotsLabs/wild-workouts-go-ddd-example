package hour

import "github.com/pkg/errors"

var (
	Available         = Availability{"available"}
	NotAvailable      = Availability{"not_available"}
	TrainingScheduled = Availability{"training_scheduled"}
)

var availabilityValues = []Availability{
	Available,
	NotAvailable,
	TrainingScheduled,
}

// Availability is enum.
//
// Using struct instead of `type Availability string` for enums allows us to ensure,
// that we have full control of what values are possible.
// With `type Availability string` you are able to create `Availability("i_can_put_anything_here")`
type Availability struct {
	a string
}

func NewAvailabilityFromString(availabilityStr string) (Availability, error) {
	for _, availability := range availabilityValues {
		if availability.String() == availabilityStr {
			return availability, nil
		}
	}
	return Availability{}, errors.Errorf("unknown '%s' availability", availabilityStr)
}

// Every type in Go have zero value. In that case it's `Availability{}`.
// It's always a good idea to check if provided value is not zero!
func (h Availability) IsZero() bool {
	return h == Availability{}
}

func (h Availability) String() string {
	return h.a
}

var (
	ErrTrainingScheduled   = errors.New("unable to modify hour, because scheduled training")
	ErrNoTrainingScheduled = errors.New("training is not scheduled")
	ErrHourNotAvailable    = errors.New("hour is not available")
)

func (h Hour) Availability() Availability {
	return h.availability
}

func (h Hour) IsAvailable() bool {
	return h.availability == Available
}

func (h Hour) HasTrainingScheduled() bool {
	return h.availability == TrainingScheduled
}

func (h *Hour) MakeNotAvailable() error {
	if h.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = NotAvailable
	return nil
}

func (h *Hour) MakeAvailable() error {
	if h.HasTrainingScheduled() {
		return ErrTrainingScheduled
	}

	h.availability = Available
	return nil
}

func (h *Hour) ScheduleTraining() error {
	if !h.IsAvailable() {
		return ErrHourNotAvailable
	}

	h.availability = TrainingScheduled
	return nil
}

func (h *Hour) CancelTraining() error {
	if !h.HasTrainingScheduled() {
		return ErrNoTrainingScheduled
	}

	h.availability = Available
	return nil
}
