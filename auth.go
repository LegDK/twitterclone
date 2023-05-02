package twitterclone

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type AuthTokenService interface {
	CreateAccessToken(ctx context.Context, user User) (string, error)
	CreateRefreshAccessToken(ctx context.Context, user User, tokenID string) (string, error)
	ParseToken(ctx context.Context, payload string) (AuthToken, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
}

type RegisterInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type AuthToken struct {
	ID  string
	Sub string
}

func (in *RegisterInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)

	in.Username = strings.TrimSpace(in.Username)
}

func (in RegisterInput) Validate() error {
	if len(in.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username not long enough, (%d) characters at least", ErrValidation, UsernameMinLength)
	}

	if !emailRegex.MatchString(in.Email) {
		return fmt.Errorf("%w: email is not valid", ErrValidation)
	}

	if len(in.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) characters at least", ErrValidation, PasswordMinLength)
	}

	if in.Password != in.ConfirmPassword {
		return fmt.Errorf("%w: password is not confirmed", ErrValidation)
	}

	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (in *LoginInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)
}

func (in LoginInput) Validate() error {
	if !emailRegex.MatchString(in.Email) {
		return fmt.Errorf("%w: email is not valid", ErrValidation)
	}

	if len(in.Password) < 1 {
		return fmt.Errorf("%w: password required", ErrValidation)
	}

	return nil
}
