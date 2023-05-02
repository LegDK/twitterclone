package domain

import (
	"context"
	"twitterclone"
)

type TweetService struct {
	TweetRepo twitterclone.TweetRepo
}

func (t TweetService) Delete(ctx context.Context, id string) error {
	currentUserID, err := twitterclone.GetUserIdFromContext(ctx)
	if err != nil {
		return twitterclone.ErrUnauthenticated
	}

	tweet, err := t.TweetRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if !tweet.CanDelete(twitterclone.User{ID: currentUserID}) {
		return twitterclone.ErrForbidden
	}

	return t.TweetRepo.Delete(ctx, id)
}

func (t TweetService) All(ctx context.Context) ([]twitterclone.Tweet, error) {
	return t.TweetRepo.All(ctx)
}

func (t TweetService) Create(ctx context.Context, input twitterclone.CreateTweetInput) (twitterclone.Tweet, error) {
	currentUserId, err := twitterclone.GetUserIdFromContext(ctx)
	if err != nil {
		return twitterclone.Tweet{}, err
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitterclone.Tweet{}, err
	}

	tweet, err := t.TweetRepo.Create(ctx, twitterclone.Tweet{UserID: currentUserId, Body: input.Body})

	if err != nil {
		return twitterclone.Tweet{}, err
	}

	return tweet, nil
}

func (t TweetService) CreateReply(ctx context.Context, input twitterclone.CreateTweetInput, parentId string) (twitterclone.Tweet, error) {
	currentUserId, err := twitterclone.GetUserIdFromContext(ctx)
	if err != nil {
		return twitterclone.Tweet{}, err
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitterclone.Tweet{}, err
	}

	if _, err := t.TweetRepo.GetById(ctx, parentId); err != nil {
		return twitterclone.Tweet{}, twitterclone.ErrNotFound
	}

	tweet, err := t.TweetRepo.Create(ctx, twitterclone.Tweet{UserID: currentUserId, Body: input.Body, ParentID: &parentId})

	if err != nil {
		return twitterclone.Tweet{}, err
	}

	return tweet, nil
}

func (t TweetService) GetById(ctx context.Context, id string) (twitterclone.Tweet, error) {
	tweet, err := t.TweetRepo.GetById(ctx, id)

	if err != nil {
		return twitterclone.Tweet{}, err
	}

	return tweet, nil
}

func NewTweetService(tr twitterclone.TweetRepo) *TweetService {
	return &TweetService{
		TweetRepo: tr,
	}
}
