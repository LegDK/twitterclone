package graph

import (
	"context"
	"twitterclone"
)

func mapTweet(tweet twitterclone.Tweet) *Tweet {
	return &Tweet{
		ID:        tweet.ID,
		UserID:    tweet.UserID,
		Body:      tweet.Body,
		CreatedAt: tweet.CreatedAt,
	}
}

func mapTweets(tweets []twitterclone.Tweet) []*Tweet {
	tt := make([]*Tweet, len(tweets))

	for i, t := range tweets {
		tt[i] = mapTweet(t)
	}

	return tt
}

func (q *queryResolver) Tweets(ctx context.Context) ([]*Tweet, error) {
	tweets, err := q.TweetService.All(ctx)
	if err != nil {
		return nil, err
	}

	return mapTweets(tweets), nil
}

func (m *mutationResolver) CreateTweet(ctx context.Context, tweetInput TweetInput) (*Tweet, error) {
	tweet, err := m.TweetService.Create(ctx, twitterclone.CreateTweetInput{
		Body: tweetInput.Body,
	})

	if err != nil {
		return nil, buildError(ctx, err)
	}

	return mapTweet(tweet), nil
}

func (m *mutationResolver) DeleteTweet(ctx context.Context, id string) (bool, error) {
	err := m.TweetService.Delete(ctx, id)

	if err != nil {
		return false, buildError(ctx, err)
	}

	return true, nil
}

func (t tweetResolver) User(ctx context.Context, obj *Tweet) (*User, error) {
	return DataloaderFor(ctx).UserByID.Load(obj.UserID)
}

func (m *mutationResolver) CreateReply(ctx context.Context, parentID string, tweetInput TweetInput) (*Tweet, error) {
	tweet, err := m.TweetService.CreateReply(ctx, twitterclone.CreateTweetInput{
		Body: tweetInput.Body,
	}, parentID)

	if err != nil {
		return nil, buildError(ctx, err)
	}

	return mapTweet(tweet), nil
}
