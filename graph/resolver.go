package graph

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	"twitterclone"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	AuthService  twitterclone.AuthService
	TweetService twitterclone.TweetService
	UserService  twitterclone.UserService
}

type queryResolver struct {
	*Resolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type tweetResolver struct {
	*Resolver
}

func (r *Resolver) Tweet() TweetResolver {
	return &tweetResolver{r}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}

func buildUnauthenticatedError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusUnauthorized,
		},
	}
}

func buildForbiddenError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildNotFoundError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, twitterclone.ErrForbidden):
		return buildForbiddenError(ctx, err)
	case errors.Is(err, twitterclone.ErrUnauthenticated):
		return buildUnauthenticatedError(ctx, err)
	case errors.Is(err, twitterclone.ErrValidation):
		return buildBadRequestError(ctx, err)
	case errors.Is(err, twitterclone.ErrNotFound):
		return buildNotFoundError(ctx, err)
	default:
		return err
	}
}
