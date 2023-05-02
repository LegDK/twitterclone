package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"runtime"
	"twitterclone/config"
)

type DB struct {
	Pool *pgxpool.Pool
	conf *config.Config
}

func New(ctx context.Context, conf *config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(conf.Database.URL)
	if err != nil {
		log.Fatalf("can't parse postgres config: %v", err)
	}

	logrusLogger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	dbConf.ConnConfig.Logger = logrusadapter.NewLogger(logrusLogger)

	conn, err := pgxpool.ConnectConfig(ctx, dbConf)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}

	db := &DB{
		Pool: conn,
		conf: conf,
	}

	return db

}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate() error {
	_, b, _, _ := runtime.Caller(0)

	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(b))
	m, err := migrate.New(migrationPath, db.conf.Database.URL)
	if err != nil {
		return fmt.Errorf("error create the migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Errorf("error migrate up: %v", err)
	}

	log.Println("migrate done")
	return nil
}

func (db *DB) Truncate(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM users`)

	if err != nil {
		return fmt.Errorf("error truncate db: %v", err)
	}

	return nil
}
