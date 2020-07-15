package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserModel struct {
	Balance     int
	DisplayName string
	Role        string
	LastIP      string
}

type db struct {
	firestoreClient *firestore.Client
}

func (d db) UsersCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("users")
}

func (d db) UserDocumentRef(userID string) *firestore.DocumentRef {
	return d.UsersCollection().Doc(userID)
}

func (d db) GetUser(ctx context.Context, userID string) (UserModel, error) {
	doc, err := d.UserDocumentRef(userID).Get(ctx)

	if err != nil && status.Code(err) != codes.NotFound {
		return UserModel{}, err
	}
	if err != nil && status.Code(err) == codes.NotFound {
		return UserModel{
			Balance: 0,
		}, nil
	}

	var user UserModel
	err = doc.DataTo(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

func (d db) UpdateBalance(ctx context.Context, userID string, amountChange int) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user UserModel

		userDoc, err := tx.Get(d.UserDocumentRef(userID))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = UserModel{
				Balance: 0,
			}
		} else {
			if err := userDoc.DataTo(&user); err != nil {
				return err
			}
		}

		user.Balance += amountChange
		if user.Balance < 0 {
			return errors.New("balance cannot be smaller than 0")
		}

		return tx.Set(userDoc.Ref, user)
	})
}

const lastIPField = "LastIP"

func (d db) UpdateLastIP(ctx context.Context, userID string, lastIP string) error {
	updates := []firestore.Update{
		{
			Path:  lastIPField,
			Value: lastIP,
		},
	}

	docRef := d.UserDocumentRef(userID)

	_, err := docRef.Update(ctx, updates)
	userNotExist := status.Code(err) == codes.NotFound

	if userNotExist {
		_, err := docRef.Set(ctx, map[string]string{lastIPField: lastIP})
		return err
	}

	return err
}
