package main

import (
	"context"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/client"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func loadFixtures(db db) {
	start := time.Now()
	logrus.Debug("Waiting for users service")

	working := client.WaitForUsersService(time.Minute * 30)
	if !working {
		logrus.Error("Users gRPC service is not up")
		return
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Users service is available")

	var attendeeUUIDs []string
	var err error
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); !mockAuth {
		for {
			attendeeUUIDs, err = createFirebaseUsers()
			if err == nil {
				logrus.Debug("Firestore users created")
				break
			}

			logrus.WithError(err).Warn("Unable to create Firestore user")
			time.Sleep(10 * time.Second)
		}
	} else {
		// ugly copy from web/src/repositories/user.js
		attendeeUUIDs = []string{"2"}
	}

	for {
		err = setAttendeeTrainingsAmount(attendeeUUIDs)
		if err == nil {
			break
		}

		logrus.WithError(err).Warn("Unable to set users credits")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Users fixtures loaded")
}

func createFirebaseUsers() ([]string, error) {
	var attendeeUUIDs []string

	var opts []option.ClientOption
	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
		opts = append(opts, option.WithCredentialsFile("/service-account-file.json"))
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		return nil, err
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	usersToCreate := []struct {
		Email       string
		DisplayName string
		Role        string
	}{
		{

			Email:       "trainer@threedots.tech",
			DisplayName: "Trainer",
			Role:        "trainer",
		},
		{
			Email:       "attendee@threedots.tech",
			DisplayName: "Mariusz Pudzianowski",
			Role:        "attendee",
		},
		{
			Email:       "attendee2@threedots.tech",
			DisplayName: "Arnold Schwarzenegger",
			Role:        "attendee",
		},
	}

	for _, user := range usersToCreate {
		userToCreate := (&auth.UserToCreate{}).
			Email(user.Email).
			Password("123456").
			DisplayName(user.DisplayName)

		createdUser, err := authClient.CreateUser(context.Background(), userToCreate)
		if err != nil && auth.IsEmailAlreadyExists(err) {
			existingUser, err := authClient.GetUserByEmail(context.Background(), user.Email)
			if err != nil {
				return nil, errors.Wrap(err, "unable to get created user")
			}
			if user.Role == "attendee" {
				attendeeUUIDs = append(attendeeUUIDs, existingUser.UID)
			}
			continue
		}
		if err != nil {
			return nil, err
		}

		err = authClient.SetCustomUserClaims(context.Background(), createdUser.UID, map[string]interface{}{
			"role": user.Role,
		})
		if err != nil {
			return nil, err
		}

		if user.Role == "attendee" {
			attendeeUUIDs = append(attendeeUUIDs, createdUser.UID)
		}
	}

	return attendeeUUIDs, nil
}

func setAttendeeTrainingsAmount(attendeeUUIDs []string) error {
	usersClient, usersClose, err := client.NewUsersClient()
	if err != nil {
		logrus.WithError(err).Error("Unable to set trainings amount")
	}
	defer usersClose()

	for _, attendeeUUID := range attendeeUUIDs {
		resp, err := usersClient.GetTrainingBalance(context.Background(), &users.GetTrainingBalanceRequest{
			UserId: attendeeUUID,
		})
		if err != nil {
			return err
		}

		if resp.Amount > 0 {
			logrus.WithFields(logrus.Fields{
				"attendee_uuid": attendeeUUID,
				"credits":       resp.Amount,
			}).Debug("Attendee have credits already")
			continue
		}

		_, err = usersClient.UpdateTrainingBalance(context.Background(), &users.UpdateTrainingBalanceRequest{
			UserId:       attendeeUUID,
			AmountChange: 20,
		})
		if err != nil {
			return err
		}

		logrus.WithField("attendee_uuid", attendeeUUID).Debug("Credits set to attendee")
	}

	return nil
}
