//go:build integration
// +build integration

package domain

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"twitterclone"
	"twitterclone/test_helpers"
)

func TestIntegrationAuthService_Register(t *testing.T) {

	validInput := twitterclone.RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register user", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(t, ctx, db)

		res, err := authService.Register(ctx, validInput)

		require.NoError(t, err)

		require.NotEmpty(t, res.User.ID)

		require.Equal(t, validInput.Email, res.User.Email)
		require.Equal(t, validInput.Username, res.User.Username)
		require.NotEqual(t, validInput.Password, res.User.Password)

	})
}
