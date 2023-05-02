package graph

import (
	"context"
	"errors"
	"twitterclone"
)

func mapAuthResponse(a twitterclone.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: a.AccessToken,
		User:        mapUser(a.User),
	}
}
func (m *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := m.AuthService.Register(ctx, twitterclone.RegisterInput{
		Username:        input.Username,
		Email:           input.Email,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, twitterclone.ErrValidation):
			return nil, buildBadRequestError(ctx, err)
		case errors.Is(err, twitterclone.ErrEmailTaken) || errors.Is(err, twitterclone.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err

		}
	}

	return mapAuthResponse(res), nil
}

func (m *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := m.AuthService.Login(ctx, twitterclone.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, twitterclone.ErrValidation):
			return nil, buildBadRequestError(ctx, err)
		case errors.Is(err, twitterclone.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err

		}
	}

	return mapAuthResponse(res), nil
}
