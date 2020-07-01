package main_test

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirestoreHourRepository(t *testing.T) {
	ctx := context.Background()
	repo := newFirebaseRepository(t, ctx)

	testCases := []struct {
		Name       string
		CreateHour func(*testing.T) *hour.Hour
	}{
		{
			Name: "available_hour",
			CreateHour: func(t *testing.T) *hour.Hour {
				return newValidAvailableHour(t, 1)
			},
		},
		{
			Name: "not_available_hour",
			CreateHour: func(t *testing.T) *hour.Hour {
				h := newValidAvailableHour(t, 2)
				require.NoError(t, h.MakeNotAvailable())

				return h
			},
		},
		{
			Name: "hour_with_training",
			CreateHour: func(t *testing.T) *hour.Hour {
				h := newValidAvailableHour(t, 3)
				require.NoError(t, h.ScheduleTraining())

				return h
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			newHour := tc.CreateHour(t)

			err := repo.UpdateHour(ctx, newHour.Time(), func(_ *hour.Hour) (*hour.Hour, error) {
				// UpdateHour provides us existing/new *hour.Hour,
				// but we are ignoring this hour and persisting result of `CreateHour`
				// we can assert this hour later in assertHourInRepository
				return newHour, nil
			})
			require.NoError(t, err)

			assertHourInRepository(ctx, t, repo, newHour)
		})
	}
}

//TestNewFirestoreHourRepository_update_existing is testing path of creating a new hour and updating this hour.
func TestNewFirestoreHourRepository_update_existing(t *testing.T) {
	ctx := context.Background()
	repo := newFirebaseRepository(t, ctx)

	testHour := newValidAvailableHour(t, 5)

	err := repo.UpdateHour(ctx, testHour.Time(), func(_ *hour.Hour) (*hour.Hour, error) {
		return testHour, nil
	})
	require.NoError(t, err)
	assertHourInRepository(ctx, t, repo, testHour)

	var expectedHour *hour.Hour
	err = repo.UpdateHour(ctx, testHour.Time(), func(h *hour.Hour) (*hour.Hour, error) {
		if err := h.ScheduleTraining(); err != nil {
			return nil, err
		}
		expectedHour = h
		return h, nil
	})
	require.NoError(t, err)
	assertHourInRepository(ctx, t, repo, expectedHour)
}

func TestNewDateDTO(t *testing.T) {

	testCases := []struct {
		Time             time.Time
		ExpectedDateTime time.Time
	}{
		{
			Time:             time.Date(3333, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedDateTime: time.Date(3333, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			// we are storing date in UTC
			// it's still 1 January 22:00 in UTC while it's midnight in +2 timezone
			Time:             time.Date(3333, 1, 2, 0, 0, 0, 0, time.FixedZone("FOO", 2*60*60)),
			ExpectedDateTime: time.Date(3333, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range testCases {
		t.Run(c.Time.String(), func(t *testing.T) {
			dateDTO := main.NewEmptyDateDTO(c.Time)
			assert.True(t, dateDTO.Date.Equal(c.ExpectedDateTime), "%s != %s", dateDTO.Date, c.ExpectedDateTime)
		})
	}
}

func newFirebaseRepository(t *testing.T, ctx context.Context) *main.FirestoreHourRepository {
	firebaseClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	require.NoError(t, err)

	repo := main.NewFirestoreHourRepository(firebaseClient)
	return repo
}

func newValidAvailableHour(t *testing.T, hourAfterMinHour int) *hour.Hour {
	hourTime := newValidHourTime(hourAfterMinHour)

	hour, err := hour.NewAvailableHour(hourTime)
	require.NoError(t, err)

	return hour
}

func newValidHourTime(hourAfterMinHour int) time.Time {
	tomorrow := time.Now().Add(time.Hour * 24)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		hour.MinUtcHour+hourAfterMinHour, 0, 0, 0,
		time.UTC,
	).Local()
}

func assertHourInRepository(ctx context.Context, t *testing.T, repo *main.FirestoreHourRepository, hour *hour.Hour) {
	hourFromRepo, err := repo.GetOrCreateHour(ctx, hour.Time())
	require.NoError(t, err)
	assert.Equal(t, hour, hourFromRepo)
}
