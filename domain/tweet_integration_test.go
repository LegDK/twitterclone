//go:build integration
// +build integration

package domain

import (
	"context"
	uuid2 "github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"twitterclone"
	"twitterclone/test_helpers"
)

func TestIntegrationTweetService_Create(t *testing.T) {
	t.Run("not auth user cannot create a tweet", func(t *testing.T) {
		ctx := context.Background()

		_, err := tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "hello",
		})

		require.ErrorIs(t, err, twitterclone.ErrUnauthenticated)
	})

	t.Run("can create a tweet", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(t, ctx, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		tweet, err := tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "hello",
		})

		require.NoError(t, err)

		require.NotEmpty(t, tweet.ID)

	})
}

func TestIntegrationTweetServiceGetAll(t *testing.T) {
	t.Run("get All", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(t, ctx, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "hello",
		})

		require.NoError(t, err)

		_, err = tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "Another one",
		})

		require.NoError(t, err)

		tweets, err := tweetService.All(ctx)

		require.NoError(t, err)

		require.NotEmpty(t, tweets)

		require.Len(t, tweets, 2)

	})
}

func TestIntegrationTweetServiceGetById(t *testing.T) {
	t.Run("get By Id", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(t, ctx, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		tweet, err := tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "hello",
		})

		require.NoError(t, err)

		require.NotEmpty(t, tweet.ID)

		newTweet, err := tweetService.GetById(ctx, tweet.ID)

		require.NoError(t, err)

		require.NotEmpty(t, newTweet.ID)

	})

	t.Run("return not found if not exists", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(t, ctx, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		tweet, err := tweetService.Create(ctx, twitterclone.CreateTweetInput{
			Body: "hello",
		})

		require.NoError(t, err)

		require.NotEmpty(t, tweet.ID)

		uuid, err := uuid2.NewV4()
		strUuid := uuid.String()

		_, err = tweetService.GetById(ctx, strUuid)

		require.ErrorIs(t, err, twitterclone.ErrNotFound)
	})
}

func TestIntegrationTweetService_Delete(t *testing.T) {
	t.Run("not auth cant delete", func(t *testing.T) {
		ctx := context.Background()

		err := tweetService.Delete(ctx, fakeUUID())
		require.ErrorIs(t, err, twitterclone.ErrUnauthenticated)

	})

	t.Run("can not delete tweet if not an owner", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(t, ctx, db)

		otherUser := test_helpers.CreateUser(ctx, t, userRepo)
		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		tweet := test_helpers.CreateTweet(ctx, t, tweetRepo, otherUser.ID)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		_, err := tweetRepo.GetById(ctx, tweet.ID)

		require.NoError(t, err)

		err = tweetService.Delete(ctx, tweet.ID)
		require.ErrorIs(t, err, twitterclone.ErrForbidden)
	})
}

func fakeUUID() string {
	uuid, _ := uuid2.NewV4()
	return uuid.String()
}
