package domain

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"twitterclone"
)

var passwordCost = bcrypt.DefaultCost

type AuthService struct {
	AuthTokenService twitterclone.AuthTokenService
	UserRepo         twitterclone.UserRepo
}

func NewAuthService(ur twitterclone.UserRepo, ats twitterclone.AuthTokenService) *AuthService {
	return &AuthService{
		UserRepo:         ur,
		AuthTokenService: ats,
	}
}

func (as *AuthService) Register(ctx context.Context, input twitterclone.RegisterInput) (twitterclone.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitterclone.AuthResponse{}, err
	}

	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, twitterclone.ErrNotFound) {
		return twitterclone.AuthResponse{}, twitterclone.ErrUsernameTaken
	}

	if _, err := as.UserRepo.GetByEmail(ctx, input.Username); !errors.Is(err, twitterclone.ErrNotFound) {
		return twitterclone.AuthResponse{}, twitterclone.ErrEmailTaken
	}

	user := twitterclone.User{
		Email:    input.Email,
		Username: input.Username,
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return twitterclone.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(hashPassword)

	user, err = as.UserRepo.Create(ctx, user)
	if err != nil {
		return twitterclone.AuthResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return twitterclone.AuthResponse{}, twitterclone.ErrGenAccessToken
	}
	return twitterclone.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil

}

func (as *AuthService) Login(ctx context.Context, input twitterclone.LoginInput) (twitterclone.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitterclone.AuthResponse{}, err
	}

	user, err := as.UserRepo.GetByEmail(ctx, input.Email)

	if err != nil {
		switch {
		case errors.Is(err, twitterclone.ErrNotFound):
			return twitterclone.AuthResponse{}, twitterclone.ErrBadCredentials
		default:
			return twitterclone.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return twitterclone.AuthResponse{}, twitterclone.ErrBadCredentials
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return twitterclone.AuthResponse{}, twitterclone.ErrGenAccessToken
	}

	return twitterclone.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}
