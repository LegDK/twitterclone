package twitterclone

import (
	"context"
	"time"
)

var (
	AccessTokenLifetime  = time.Minute * 15
	RefreshTokenLifeTime = time.Hour * 24 * 7
)

type RefreshToken struct {
	ID         string
	Name       string
	UserID     string
	LastUsedAt time.Time
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

type CreateRefreshTokenParams struct {
	Sub  string
	Name string
}
type RefreshTokenRepo interface {
	Create(ctx context.Context, params CreateRefreshTokenParams) (RefreshToken, error)
	GetById(ctx context.Context, id string) (RefreshToken, error)
}
