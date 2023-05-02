package twitterclone

import (
	"context"
	"errors"
)

type contextKey string

var (
	contextAuthIDKey contextKey = "currentUserId"
)

func GetUserIdFromContext(ctx context.Context) (string, error) {

	if ctx.Value(contextAuthIDKey) == nil {
		return "", ErrUnauthenticated
	}

	userID, ok := ctx.Value(contextAuthIDKey).(string)
	if !ok {
		return "", errors.New("no user id in context")
	}
	return userID, nil
}

func PutUserIdIntoContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextAuthIDKey, userID)
}
