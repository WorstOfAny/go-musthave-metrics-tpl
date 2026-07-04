package repository

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/migrations"
	models "github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
)

type Config struct {
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	RestoreStorage  bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

type Repository interface {
	Set(context.Context, models.Metrics) error
	BulkSet(context.Context, []models.Metrics) error
	Get(context.Context, string) (models.Metrics, error)
	Remove(context.Context, string) error
	All(context.Context) (iter.Seq[models.Metrics], error)
	Ping(context.Context) error
}

type repositoryError string

func (re repositoryError) Error() string {
	return string(re)
}

const ErrNotFound = repositoryError("metric not found")

func NewRepository(ctx context.Context, cfg Config) (Repository, error) {

	switch {
	case cfg.DatabaseDSN != "":
		sourceDriver, err := iofs.New(migrations.MigrationsFS, "migrations")
		m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed initialize migrations: %w", err)
		}
		err = m.Up()

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}

		db, err := pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize db connections pool: %w", err)
		}

		go func() {
			<-ctx.Done()
			db.Close()
		}()

		return NewDBDecorator(db), nil
	case cfg.FileStoragePath != "":
		file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		go func() {
			<-ctx.Done()
			file.Close()
		}()

		repo, err := NewStorage(file, cfg.RestoreStorage)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mem storage: %w", err)
		}
		go repo.WriteToFile(ctx, time.Duration(cfg.StoreInterval)*time.Second)
		return repo, nil
	default:
		repo, err := NewStorage(nil, false)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mem storage: %w", err)
		}
		return repo, nil
	}

}
