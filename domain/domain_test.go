package domain

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
	"twitterclone"
	"twitterclone/config"
	"twitterclone/jwt"
	"twitterclone/postgres"
)

var (
	conf             *config.Config
	db               *postgres.DB
	authService      twitterclone.AuthService
	userRepo         twitterclone.UserRepo
	tweetService     twitterclone.TweetService
	tweetRepo        twitterclone.TweetRepo
	authTokenService twitterclone.AuthTokenService
	contextAuthIDKey = "currentUserId"
)

func TestMain(t *testing.M) {
	passwordCost = bcrypt.MinCost

	config.LoadEnv(".env.test")

	conf = config.New()

	fmt.Println(conf)
	db = postgres.New(context.Background(), conf)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userRepo = postgres.NewUserRepo(db)
	authTokenService = jwt.NewTokenService(conf)
	authService = NewAuthService(userRepo, authTokenService)
	tweetRepo = postgres.NewTweetRepo(db)
	tweetService = NewTweetService(tweetRepo)
	os.Exit(t.Run())
}
