package main

import (
	"net/http"
	"twitterclone"
)

func authMiddleware(authTokenService twitterclone.AuthTokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()

			token, err := authTokenService.ParseTokenFromRequest(ctx, request)

			if err != nil {
				next.ServeHTTP(writer, request)
				return
			}

			ctx = twitterclone.PutUserIdIntoContext(ctx, token.Sub)

			next.ServeHTTP(writer, request.WithContext(ctx))

		})
	}
}
