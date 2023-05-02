package jwt

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	_ "github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"time"
	"twitterclone"
	"twitterclone/config"
)

var signatureType = jwa.HS256

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (twitterclone.AuthToken, error) {
	token, err := jwtGo.ParseRequest(
		r,
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)

	if err != nil {
		return twitterclone.AuthToken{}, twitterclone.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func buildToken(token jwtGo.Token) twitterclone.AuthToken {
	return twitterclone.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (twitterclone.AuthToken, error) {
	token, err := jwtGo.Parse(
		[]byte(payload),
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)

	if err != nil {
		return twitterclone.AuthToken{}, twitterclone.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user twitterclone.User) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, twitterclone.AccessTokenLifetime, ts.Conf); err != nil {
		return "", err
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (ts *TokenService) CreateRefreshAccessToken(ctx context.Context, user twitterclone.User, tokenID string) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, twitterclone.RefreshTokenLifeTime, ts.Conf); err != nil {
		return "", err
	}

	if err := t.Set(jwtGo.JwtIDKey, tokenID); err != nil {
		return "", fmt.Errorf("error set jwt id: %v", err)
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func setDefaultToken(t jwtGo.Token, user twitterclone.User, lifetime time.Duration, conf *config.Config) error {
	if err := t.Set(jwtGo.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("error set jwt sub: %v", err)
	}

	if err := t.Set(jwtGo.IssuerKey, conf.JWT.Issuer); err != nil {
		return fmt.Errorf("error set jwt issuer: %v", err)
	}

	if err := t.Set(jwtGo.IssuedAtKey, time.Now().Unix()); err != nil {
		return fmt.Errorf("error set jwt issuered at key: %v", err)
	}

	if err := t.Set(jwtGo.ExpirationKey, time.Now().Add(lifetime).Unix()); err != nil {
		return fmt.Errorf("error set jwt expiration : %v", err)
	}

	return nil
}
