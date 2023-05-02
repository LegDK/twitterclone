package twitterclone

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	TweetMinLength = 2
	TweetMaxLength = 255
)

type CreateTweetInput struct {
	Body string
}

func (in *CreateTweetInput) Sanitize() {
	in.Body = strings.TrimSpace(in.Body)
}

func (in CreateTweetInput) Validate() error {
	if len(in.Body) < TweetMinLength || len(in.Body) > TweetMaxLength {
		return fmt.Errorf("%w: Body not long or to long, range: %d to %d", ErrValidation, TweetMinLength, TweetMaxLength)
	}
	return nil
}

type Tweet struct {
	ID        string
	Body      string
	UserID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Tweet) CanDelete(user User) bool {
	return t.UserID == user.ID
}

type TweetService interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, input CreateTweetInput) (Tweet, error)
	CreateReply(ctx context.Context, input CreateTweetInput, parentID string) (Tweet, error)
	GetById(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}

type TweetRepo interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, input Tweet) (Tweet, error)
	GetById(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}
