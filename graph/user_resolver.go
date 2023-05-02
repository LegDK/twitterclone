package graph

import (
	"context"
	"twitterclone"
)

func (q *queryResolver) Me(ctx context.Context) (*User, error) {
	userID, err := twitterclone.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, twitterclone.ErrUnauthenticated
	}

	return mapUser(twitterclone.User{
		ID: userID,
	}), nil
}

func mapUser(u twitterclone.User) *User {
	return &User{
		Username:  u.Username,
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
