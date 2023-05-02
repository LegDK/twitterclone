package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"twitterclone"
)

type TweetRepo struct {
	DB *DB
}

func (t TweetRepo) Delete(ctx context.Context, id string) error {
	tx, err := t.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	err = deleteTweet(ctx, tx, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (t TweetRepo) All(ctx context.Context) ([]twitterclone.Tweet, error) {
	query := `SELECT * FROM tweets ORDER BY created_at DESC;`

	var tweets []twitterclone.Tweet

	if err := pgxscan.Select(ctx, t.DB.Pool, &tweets, query); err != nil {
		if pgxscan.NotFound(err) {
			return []twitterclone.Tweet{}, twitterclone.ErrNotFound
		}

		return []twitterclone.Tweet{}, fmt.Errorf("error select: %v", err)
	}

	return tweets, nil
}

func (t TweetRepo) Create(ctx context.Context, tweet twitterclone.Tweet) (twitterclone.Tweet, error) {
	tx, err := t.DB.Pool.Begin(ctx)
	if err != nil {
		return twitterclone.Tweet{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	tweet, err = createTweet(ctx, tx, tweet)
	if err != nil {
		return twitterclone.Tweet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitterclone.Tweet{}, err
	}

	return tweet, nil
}

func (t TweetRepo) GetById(ctx context.Context, id string) (twitterclone.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id = $1;`

	tweet := twitterclone.Tweet{}

	if err := pgxscan.Get(ctx, t.DB.Pool, &tweet, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitterclone.Tweet{}, twitterclone.ErrNotFound
		}

		return twitterclone.Tweet{}, fmt.Errorf("error select: %v", err)
	}

	return tweet, nil
}

func NewTweetRepo(db *DB) *TweetRepo {
	return &TweetRepo{
		DB: db,
	}
}

func createTweet(ctx context.Context, tx pgx.Tx, tweet twitterclone.Tweet) (twitterclone.Tweet, error) {
	query := `INSERT INTO tweets(body, user_id, parent_id) VALUES ($1, $2, $3) RETURNING *;`

	t := twitterclone.Tweet{}

	if err := pgxscan.Get(ctx, tx, &t, query, tweet.Body, tweet.UserID, tweet.ParentID); err != nil {
		return twitterclone.Tweet{}, fmt.Errorf("error insert: %v", err)
	}

	return t, nil
}

func deleteTweet(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM tweets WHERE id = $1;`

	if _, err := tx.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("error insert: %v", err)
	}

	return nil
}
