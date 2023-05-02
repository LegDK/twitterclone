package domain

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"twitterclone"
	"twitterclone/mocks"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twitterclone.RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twitterclone.User{
			ID:       "123",
			Username: validInput.Username,
			Email:    validInput.Email,
		}, nil)
		authTokenService := &mocks.AuthTokenService{}
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("123", nil)
		service := NewAuthService(userRepo, authTokenService)
		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)

		userRepo.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitterclone.User{}, nil)

		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitterclone.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{}, nil)

		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitterclone.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)
		userRepo.On("Create", mock.Anything, mock.Anything).Return(twitterclone.User{}, errors.New("something"))

		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		validInput.Email = "fuck"
		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {

	validInput := twitterclone.LoginInput{
		Email:    "bob@gmail.com",
		Password: "password",
	}

	t.Run("can login", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(validInput.Password), bcrypt.MinCost)
		require.NoError(t, err)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{
			Email:    validInput.Email,
			Password: string(hashPassword),
		}, nil)
		authTokenService := &mocks.AuthTokenService{}
		authTokenService.On("CreateAccessToken", mock.Anything, mock.Anything).Return("123", nil)
		service := NewAuthService(userRepo, authTokenService)
		_, err = service.Login(ctx, validInput)
		require.NoError(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(validInput.Password), bcrypt.MinCost)
		require.NoError(t, err)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{
			Email:    validInput.Email,
			Password: string(hashPassword),
		}, nil)

		validInput.Password = "another password"
		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err = service.Login(ctx, validInput)
		require.ErrorIs(t, err, twitterclone.ErrBadCredentials)

		userRepo.AssertExpectations(t)
	})

	t.Run("email not found", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(twitterclone.User{}, twitterclone.ErrNotFound)

		authTokenService := &mocks.AuthTokenService{}
		service := NewAuthService(userRepo, authTokenService)
		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, twitterclone.ErrBadCredentials)

		userRepo.AssertExpectations(t)
	})
}
