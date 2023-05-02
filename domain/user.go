package domain

import (
	"context"
	"twitterclone"
)

type UserService struct {
	UserRepo twitterclone.UserRepo
}

func (u UserService) GetById(ctx context.Context, id string) (twitterclone.User, error) {
	return u.UserRepo.GetById(ctx, id)
}

func NewUserSerivce(ur twitterclone.UserRepo) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}
