package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
	"twitterclone/config"
	"twitterclone/domain"
	"twitterclone/graph"
	"twitterclone/jwt"
	"twitterclone/postgres"
)

func main() {
	ctx := context.Background()

	config.LoadEnv(".env")

	conf := config.New()

	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepo(db)
	tr := postgres.NewTweetRepo(db)
	authTokenService := jwt.NewTokenService(conf)
	authService := domain.NewAuthService(userRepo, authTokenService)
	tweetService := domain.NewTweetService(tr)
	userService := domain.NewUserSerivce(userRepo)

	r := chi.NewRouter()
	r.Use(authMiddleware(authTokenService))
	r.Use(graph.DataloaderMiddleware(&graph.Repos{
		UserRepo: userRepo,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Handle("/", playground.Handler("Twitter clone", "/query"))
	r.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService:  authService,
					TweetService: tweetService,
					UserService:  userService,
				},
			},
		),
	))
	log.Fatal(http.ListenAndServe(":3000", r))
	log.Println("Server running")
}
