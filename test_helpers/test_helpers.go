package test_helpers

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"twitterclone"
	"twitterclone/faker"
	"twitterclone/postgres"
)

func TeardownDB(t *testing.T, ctx context.Context, db *postgres.DB) {

	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, userRepo twitterclone.UserRepo) twitterclone.User {
	t.Helper()

	user, err := userRepo.Create(ctx, twitterclone.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password,
	})
	require.NoError(t, err)

	return user
}

func CreateTweet(ctx context.Context, t *testing.T, repo twitterclone.TweetRepo, forUser string) twitterclone.Tweet {
	t.Helper()

	tweet, err := repo.Create(ctx, twitterclone.Tweet{
		Body:   "fuck",
		UserID: forUser,
	})

	require.NoError(t, err)

	return tweet
}

func LoginUser(ctx context.Context, t *testing.T, user twitterclone.User) context.Context {
	t.Helper()

	return twitterclone.PutUserIdIntoContext(ctx, user.ID)
}
