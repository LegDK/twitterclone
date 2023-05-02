package twitterclone

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserIdFromContext(t *testing.T) {

	t.Run("get user context", func(t *testing.T) {

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextAuthIDKey, "123")
		userID, err := GetUserIdFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})

}

func TestPutUserIdIntoContext(t *testing.T) {
	t.Run("get user context", func(t *testing.T) {

		ctx := context.Background()
		ctx = PutUserIdIntoContext(ctx, "123")

		userID, err := GetUserIdFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})
}
