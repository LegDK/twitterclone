package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"twitterclone"
)

type UserRepo struct {
	DB *DB
}

func (ur *UserRepo) GetByIds(ctx context.Context, ids []string) ([]twitterclone.User, error) {
	query := `SELECT * FROM users WHERE id = ANY($1);`

	var users []twitterclone.User

	if err := pgxscan.Select(ctx, ur.DB.Pool, &users, query, ids); err != nil {
		if pgxscan.NotFound(err) {
			return []twitterclone.User{}, twitterclone.ErrNotFound
		}

		return []twitterclone.User{}, fmt.Errorf("error select: %v", err)
	}

	return users, nil
}

func (ur *UserRepo) GetById(ctx context.Context, id string) (twitterclone.User, error) {
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1;`

	u := twitterclone.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitterclone.User{}, twitterclone.ErrNotFound
		}

		return twitterclone.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (ur *UserRepo) Create(ctx context.Context, user twitterclone.User) (twitterclone.User, error) {

	tx, err := ur.DB.Pool.Begin(ctx)
	if err != nil {
		return twitterclone.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return twitterclone.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitterclone.User{}, err
	}

	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user twitterclone.User) (twitterclone.User, error) {
	query := `INSERT INTO users(email, username, password) VALUES ($1, $2, $3) RETURNING *;`

	u := twitterclone.User{}

	if err := pgxscan.Get(ctx, tx, &u, query, user.Email, user.Username, user.Password); err != nil {
		return twitterclone.User{}, fmt.Errorf("error insert: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (twitterclone.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1;`

	u := twitterclone.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return twitterclone.User{}, twitterclone.ErrNotFound
		}

		return twitterclone.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (twitterclone.User, error) {

	query := `SELECT * FROM users WHERE email = $1 LIMIT 1;`

	u := twitterclone.User{}

	if err := pgxscan.Get(ctx, ur.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return twitterclone.User{}, twitterclone.ErrNotFound
		}

		return twitterclone.User{}, fmt.Errorf("error select: %v", err)
	}

	return u, nil
}
