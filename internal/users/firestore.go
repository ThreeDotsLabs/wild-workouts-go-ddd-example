package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type db struct {
	firestoreClient *firestore.Client
}

func (d db) UsersCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("users")
}

func (d db) UserDocumentRef(userID string) *firestore.DocumentRef {
	return d.UsersCollection().Doc(userID)
}

func (d db) GetUser(ctx context.Context, userID string) (User, error) {
	doc, err := d.UserDocumentRef(userID).Get(ctx)

	if err != nil && status.Code(err) != codes.NotFound {
		return User{}, err
	}
	if err != nil && status.Code(err) == codes.NotFound {
		return User{
			Balance: 0,
		}, nil
	}

	var user User
	err = doc.DataTo(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (d db) UpdateBalance(ctx context.Context, userID string, amountChange int) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user User

		userDoc, err := tx.Get(d.UserDocumentRef(userID))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = User{
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
