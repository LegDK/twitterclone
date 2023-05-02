package jwt

import (
	"context"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
	"twitterclone"
	"twitterclone/config"
)

var (
	conf         *config.Config
	tokenService *TokenService
	now          func() time.Time
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")
	conf = config.New()

	tokenService = NewTokenService(conf)
	os.Exit(m.Run())
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("can create a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitterclone.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithVerify(signatureType, []byte(conf.JWT.Secret)),
			jwtGo.WithIssuer(conf.JWT.Issuer),
		)

		require.NoError(t, err)

		require.Equal(t, "123", tok.Subject())
		require.Equal(t, time.Now().Add(twitterclone.AccessTokenLifetime).Unix(), tok.Expiration().Unix())
	})
}

func TestTokenService_CreateRefreshAccessToken(t *testing.T) {
	t.Run("can create a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitterclone.User{
			ID: "123",
		}

		token, err := tokenService.CreateRefreshAccessToken(ctx, user, "456")
		require.NoError(t, err)

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithVerify(signatureType, []byte(conf.JWT.Secret)),
			jwtGo.WithIssuer(conf.JWT.Issuer),
		)

		require.NoError(t, err)

		require.Equal(t, "123", tok.Subject())
		require.Equal(t, "456", tok.JwtID())
		require.Equal(t, time.Now().Add(twitterclone.RefreshTokenLifeTime).Unix(), tok.Expiration().Unix())

		teardownTimeNow(t)
	})
}

func TestTokenService_ParseToken(t *testing.T) {
	t.Run("can parse a valid token", func(t *testing.T) {
		ctx := context.Background()
		user := twitterclone.User{
			ID: "123",
		}

		token, err := tokenService.CreateRefreshAccessToken(ctx, user, "456")
		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, "123", tok.Sub)
	})
}

func teardownTimeNow(t *testing.T) {
	t.Helper()

	now = func() time.Time {
		return time.Now()
	}
}
